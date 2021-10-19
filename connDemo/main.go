package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/lucas-clemente/quic-go"
)

const DEFAULT_CERT = "cert.pem"
const DEFAULT_KEY = "key.pem"

const USAGE = `Usage:
-t          tcp mode
-u          udp mode
-q          quic mode
-s [addr]   serve mode, will listen to the specified address with the specified protocol. for example: '-s 0.0.0.0:8080'
-c [addr]   client mode, will connect to the address by the specified protocol, and will be overrided by '-s'. for example: '-t 127.0.0.1:8080'
-tls        enable tls

for example:
go run main.go -t -s 127.0.0.1:8081 -tls # create a tcp server with tls
go run main.go -t -s 127.0.0.1:8081      # create a tcp server
go run main.go -u -s 127.0.0.1:8081      # create a udp server
go run main.go -q -s 127.0.0.1:8081      # create a quic server

go run main.go -t -c 127.0.0.1:8081 -tls # create a tcp client to connect to a tls server
go run main.go -t -c 127.0.0.1:8081      # create a tcp client
go run main.go -t -c 127.0.0.1:8081      # create a udp client
go run main.go -q -c 127.0.0.1:8081      # create a quic client
`

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func main() {
	var protocal string
	var serve, client bool
	var isTLS bool
	var addr string
	if len(os.Args) < 2 {
		fmt.Println(USAGE)
		return
	}
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-t":
			protocal = "tcp"
		case "-u":
			protocal = "tcp"
		case "-q":
			protocal = "quic"
		case "-s":
			serve = true
			addr = os.Args[i+1]
			i++
		case "-c":
			client = true
			addr = os.Args[i+1]
			i++
		case "-tls":
			isTLS = true
		default:
			fmt.Println(USAGE)
			os.Exit(1)
		}
	}
	// fmt.Println(tcpMode, udpMode, serve, client, isTLS, addr)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go handleSignals(cancel)

	if serve {
		switch protocal {
		case "tcp":
			if isTLS {
				tlsServe(ctx, addr)
				return
			}
			tcpServe(ctx, addr)
		case "udp":
			if isTLS {
				fmt.Println("udp not support tls. you can try quic protocol\n", USAGE)
				os.Exit(0)
			}
			udpServe(ctx, addr)
		case "quic":
			quicServe(ctx, addr)
		default:
			fmt.Println(USAGE)
			os.Exit(1)
		}
		return
	}

	if client {
		switch protocal {
		case "tcp":
			if isTLS {
				tlsClient(ctx, addr)
				return
			}
			tcpClient(ctx, addr)
		case "udp":
			if isTLS {
				fmt.Println("udp not support tls. you can try quic protocol\n", USAGE)
				os.Exit(0)
			}
			udpClient(ctx, addr)
		case "quic":
			quicClient(ctx, addr)
		default:
			fmt.Println(USAGE)
			os.Exit(1)
		}
	}
}

func isFileExist(name string) (bool, error) {
	info, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return !info.IsDir(), nil
}

func getSSLKeyFile() (f *os.File, err error) {
	filename := "sslkeylog.log"
	exist, err := isFileExist(filename)
	if err != nil {
		err = fmt.Errorf("can not write file: %v", err)
		return nil, err
	}
	if exist {
		f, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644) //打开文件
		if err != nil {
			err = fmt.Errorf("can not open file: %w", err)
			return nil, err
		}
	} else {
		f, err = os.Create(filename) //创建文件
		if err != nil {
			err = fmt.Errorf("can not create file: %v", err)
			return nil, err
		}
	}
	return
}

func genTLS(altIPs, altDNSs []string, enableKeyLogFile, writeToDisk bool) (*tls.Config, error) {
	ips := []net.IP{}
	addlocalhost := true
	for _, ip := range altIPs {
		newIP := net.ParseIP(ip)
		if newIP == nil {
			return nil, fmt.Errorf("%q not a IP address", ip)
		}
		ips = append(ips, newIP)
		if ip == "127.0.0.1" {
			addlocalhost = false
		}
	}
	if addlocalhost {
		ips = append(ips, net.ParseIP("127.0.0.1"))
	}

	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, max)
	subject := pkix.Name{
		Country:      []string{"CN"},
		Province:     []string{"BeiJing"},
		Organization: []string{"test"},
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour), // 1 year
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  ips,
		DNSNames:     altDNSs,
	}

	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("gen rsa key failed: %w", err)
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk)
	if err != nil {
		return nil, fmt.Errorf("CreateCertificate failed: %w", err)
	}
	certPem := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyPem := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	tlsCert, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		return nil, fmt.Errorf("create x509 pair: %w", err)
	}

	if writeToDisk {
		tmp := [2][]byte{certPem, keyPem}
		for i, v := range [2]string{DEFAULT_CERT, DEFAULT_KEY} {
			if err = ioutil.WriteFile(v, tmp[i], 0644); err != nil {
				return nil, fmt.Errorf("write %q: %w", v, err)
			}
		}
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic"}, // is required for quic
	}
	if !enableKeyLogFile {
		return tlsConfig, nil
	}

	f, err := getSSLKeyFile()
	if err != nil {
		err = fmt.Errorf("can not get ssl key file: %v", err)
		return nil, err
	}
	tlsConfig.KeyLogWriter = f

	return tlsConfig, err
}

func getIPPort(addr string) (ip string, port int, err error) {
	addrs := strings.Split(addr, ":")
	switch len(addrs) {
	case 2:
		ip = addrs[0]
		port, err = strconv.Atoi(addrs[1])
	case 1:
		err = fmt.Errorf("%q without a port", addr)
	default:
		err = fmt.Errorf("%q not a ipv4 address", addr)
	}
	return
}

type MyConn interface {
	// Read reads data from the connection.
	// Read can be made to time out and return an error after a fixed
	// time limit; see SetDeadline and SetReadDeadline.
	Read(b []byte) (n int, err error)

	// Write writes data to the connection.
	// Write can be made to time out and return an error after a fixed
	// time limit; see SetDeadline and SetWriteDeadline.
	Write(b []byte) (n int, err error)

	// Close closes the connection.
	// Any blocked Read or Write operations will be unblocked and return errors.
	Close() error
}

func handleConn(conn MyConn, remoteAddr string) {
	defer conn.Close()
	log.Println("Receive Connect Request From", remoteAddr)
	buffer := make([]byte, 1024)
	for {
		len, err := conn.Read(buffer)
		if err != nil {
			log.Println("can not read from conn:", err)
			break
		}
		fmt.Printf("Receive from[%s]: %s\n", remoteAddr, string(buffer[:len]))
		// _, err = conn.Write([]byte("thanks"))
		_, err = conn.Write(append([]byte("server got: "), buffer[:len]...))
		if err != nil {
			log.Println("can not write to conn:", err)
			break
		}
	}
	fmt.Println("Client " + remoteAddr + " Connection Closed.....")
}

func handleTCPConn(conn net.Conn) {
	handleConn(conn, conn.RemoteAddr().String())
}

func tlsServe(ctx context.Context, addr string) {
	ip, _, err := getIPPort(addr)
	if err != nil {
		log.Fatal(err)
	}
	tlsConfig, err := genTLS([]string{ip}, nil, true, false)
	if err != nil {
		log.Fatal(err)
	}
	listener, err := tls.Listen("tcp", addr, tlsConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	log.Printf("listenning %s on tcp with tls", addr)
	for {
		select {
		case <-ctx.Done():
			log.Printf("stop listenning %s on TCP with tls. good bye!", addr)
			return
		default:
		}

		conn, err := listener.Accept()
		if err != nil {
			log.Println("connect err:", err)
			continue
		}
		go handleTCPConn(conn)
	}
}

func tcpServe(ctx context.Context, addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	log.Printf("listenning %s on tcp", addr)
	for {
		select {
		case <-ctx.Done():
			log.Printf("stop listenning %s on TCP. good bye!", addr)
			return
		default:
		}

		conn, err := listener.Accept()
		if err != nil {
			log.Println("connect err:", err)
			continue
		}
		go handleTCPConn(conn)
	}
}

func udpServe(ctx context.Context, addr string) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	listener, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	log.Printf("listenning %s on udp", addr)
	buffer := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			log.Printf("stop listenning %s on UDP. good bye!", addr)
			return
		default:
		}

		n, remoteAddr, err := listener.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("error during read: %s", err)
		}
		log.Printf("<%s> %s\n", remoteAddr.String(), buffer[:n])
		_, err = listener.WriteToUDP(append([]byte("server got: "), buffer...), remoteAddr)
		if err != nil {
			log.Printf("error during response: %s", err)
		}
	}
}

func quicServe(ctx context.Context, addr string) {
	ip, _, err := getIPPort(addr)
	if err != nil {
		log.Fatal(err)
	}
	tlsConfig, err := genTLS([]string{ip}, nil, true, false)
	if err != nil {
		log.Fatal(err)
	}
	listener, err := quic.ListenAddr(addr, tlsConfig, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	log.Printf("listenning %s on quic", addr)
	for {
		select {
		case <-ctx.Done():
			log.Printf("stop listenning %s on quic. good bye!", addr)
			return
		default:
		}
		sess, err := listener.Accept(ctx)
		if err != nil {
			log.Println("connect err:", err)
			continue
		}
		go func(sess quic.Session) {
			str, err := sess.AcceptStream(context.Background())
			if err != nil {
				log.Printf("failed to connect: %v", err)
				return
			}
			handleConn(str, sess.RemoteAddr().String())
		}(sess)
	}
}

func openClient(ctx context.Context, conn MyConn, remoteAddr, protocol string, withTLS bool) {
	ticker := time.NewTicker(1 * time.Millisecond * 3000)
	reader := make(chan []byte)
	cancelctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		buf := make([]byte, 1024)
		for {
			select {
			case <-cancelctx.Done():
				return
			default:
			}
			_, err := conn.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Println(err)
				}
				return
			}
			reader <- buf
		}
	}()
	for {
		select {
		case <-ctx.Done():
			if withTLS {
				log.Printf("stop connecting %s on %s. good bye!", remoteAddr, protocol)
			} else {
				log.Printf("stop connecting %s on %s with tls. good bye!", remoteAddr, protocol)
			}
			return
		case <-ticker.C:
			// _, err = io.WriteString(conn, "hello, now is "+time.Now().String())
			_, err := conn.Write([]byte("hello, now is " + time.Now().String()))
			if err != nil {
				log.Fatalln(err)
			}
		case buf := <-reader:
			log.Println("Receive From Server:", string(buf[:]))
		}
	}
}

func clientConn(ctx context.Context, conn net.Conn, protocol string, withTLS bool) {
	openClient(ctx, conn, conn.RemoteAddr().String(), protocol, withTLS)
}

func tlsClient(ctx context.Context, addr string) {
	conn, err := tls.Dial("tcp", addr, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		log.Fatal("dial failed: ", err)
	}
	defer conn.Close()
	log.Println("Client Connect To ", conn.RemoteAddr())
	status := conn.ConnectionState()
	log.Printf("connect status %#v\n", status)
	clientConn(ctx, conn, "TCP", true)
}

func tcpClient(ctx context.Context, addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal("dial failed:", err)
	}
	defer conn.Close()
	clientConn(ctx, conn, "TCP", false)
}

func udpClient(ctx context.Context, addr string) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal("dial failed:", err)
	}
	defer conn.Close()
	clientConn(ctx, conn, "UDP", false)
}

func quicClient(ctx context.Context, addr string) {
	f, err := getSSLKeyFile()
	if err != nil {
		log.Fatal(err)
	}
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic"},
		KeyLogWriter:       f,
	}
	sess, err := quic.DialAddr(addr, tlsConf, nil)
	if err != nil {
		log.Fatal("dial failed: ", err)
	}
	log.Println("Client Connect To ", sess.RemoteAddr())
	status := sess.ConnectionState()
	log.Printf("connect status %#v\n", status)
	str, err := sess.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatal("connect stream failed: ", err)
	}
	openClient(ctx, str, sess.RemoteAddr().String(), "QUIC", true)
}

func handleSignals(cancel context.CancelFunc) {
	sigc := make(chan os.Signal, 2)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)

	for sig := range sigc {
		switch sig {
		case syscall.SIGINT:
			cancel()
			fmt.Println("good bye!")
			os.Exit(1)
		case syscall.SIGTERM:
			cancel()
			return
		}
	}
}
