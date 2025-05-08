package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

type fedControlClusterClient struct {
	count      int32
	cli        client.Client
	kubeconfig string
}

var (
	eg                       errgroup.Group
	lock                     sync.Mutex
	fedControlClusterClients = map[string]*fedControlClusterClient{}
	scheme                   = runtime.NewScheme()
)

var kubeconfig = `apiVersion: v1
clusters:
- cluster:
    #certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM1ekNDQWMrZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRJek1ESXhNREEyTVRnME0xb1hEVE16TURJd056QTJNVGcwTTFvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTW05ClNCYjlGdmQ0eWx2UmVjMWZzWEVtR25TR0FncmxkaVFvU2d1aXhmNXRRUVo1VDQwYUtPbXNEQmc4Nk4wdEM4UjYKZ3lvRzVKb2ZMa0NZT3VpRDRCcEd4UmdvbXlJd3BTNDMvbzlDNHI2R0Q3MWJ6VFliMWNzWlZIbmJWdU1OS3RYTwpFUjBHU080Nk9xaUFCazhkQk5qUjlPVFVyUDhTU2VyeDVUREtxSmVHMnNtdDFWSE5vODNlczQ1SzArTitNb0M0CmdVNHdHNGtOMEtFYUo1eUozd3p0YnlCR3d6bFZKMkM2RjRFL1p3cGFibG0zeW9HRXNvMG1BbVFIdWFjb3htM3gKdnpaMEh4M3ZBOTBnME9IQjBpVytxS1VjSk8vN2xDMVpVNVlvSzgwbU5tbndodnozK01PbjFPekJTMlRneVUxQgpJQlZjT1k5TnVmRStXdVlrMWtFQ0F3RUFBYU5DTUVBd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0hRWURWUjBPQkJZRUZGaTZPcTV2Yit4dVJBc29zZlhoRVI4L2ZHZG9NQTBHQ1NxR1NJYjMKRFFFQkN3VUFBNElCQVFBY3VOaXoyQ0tTbUxrR0UrWkNvUC8vU0pWV2ZqeElGMWRxSmdlcXhoUmpDMTVxRHJaZAp2ZWpONFdyMlhCdS9JY3ZGM0tqU1RMbGx1Q2puUXBCSWlJNHJ6dkd5UHRObHlnc2pJdDA2N3J6NkVLWTIyU05RCldrT2xmdUxzWFZjV3llZkM3ZXNzM1puSUsxSHpaMzdhYjRodWc2SDNtUmZ2aWlyRGp2d0kvU0QzK3RPamkwZm8KMk1EV1lmZStlVldyNzBkOTB3bzJuckdhQmMrOXpkdjJHdnZYSDk1RHZlT3YyYkE0WEduQnpvQVRWWFY1c2MxVwpURFRYSFVxcFdQazh3NWowNExkcFczNi9BUFdPWFZiS254TzRtbkJXUFZpZ1pnWVZMdHhYMHl5YWdLdmwwQ1ZoClIxMDc5SzAydkwvNnYyeWYxLzRvbnplbC8ycjRmTm55UFFvUwotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
    insecure-skip-tls-verify: true
    server: https://10.37.13.21:34783
  name: kind-dcp-test
contexts:
- context:
    cluster: kind-dcp-test
    user: kind-dcp-test
  name: kind-dcp-test
current-context: kind-dcp-test
kind: Config
preferences: {}
users:
- name: kind-dcp-test
  user:
    client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURFekNDQWZ1Z0F3SUJBZ0lJWTY2dzJBTFF5cTR3RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB5TXpBeU1UQXdOakU0TkROYUZ3MHlOREF5TVRBd05qRTRORFZhTURReApGekFWQmdOVkJBb1REbk41YzNSbGJUcHRZWE4wWlhKek1Sa3dGd1lEVlFRREV4QnJkV0psY201bGRHVnpMV0ZrCmJXbHVNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQTFoclRSdy9wc2hVMjR0VFYKSmU2T0FPQUpIL2NiNUdHWUtwRVYya29TNHR3RDFyUUNPTU4yYXVoMVdzK05PMDVwZXZwQnE2NDc0RW4vSXJMQwpSSnFNT010YmZYS1YxdHVydzB0WDlBd3dsb2ZxWjRGWUwzaUpTVkQ1TU9ISXZrWkdjV3ppOHh4SGcxdzllUzZLCm9lcEZ4dHZWSUIwU2U4ZVFpQUN4Y1F2L2JWTmlYRjQ1M0t6UFBNZjRvWCs1MG1oL01oa1J3dGQrdUt6Ym1LaE4KbzhPQUY0T2lGRHBTbXlKdGNlVXBBR0JKSTR6TDRoRVJOUmVBTzVYL2lRZFdFNVpiN2lYekFVd01SKysvcGRtNAowVnQ0TmZ2TCs2Mll0a3ZpVnQ3Mnl3Z05JZmFqZ2JXUjlWaEdHYWxyYjZmY0xHa0kwekk5THdkZVltVHdyaHd5ClY2Mjh5d0lEQVFBQm8wZ3dSakFPQmdOVkhROEJBZjhFQkFNQ0JhQXdFd1lEVlIwbEJBd3dDZ1lJS3dZQkJRVUgKQXdJd0h3WURWUjBqQkJnd0ZvQVVXTG82cm05djdHNUVDeWl4OWVFUkh6OThaMmd3RFFZSktvWklodmNOQVFFTApCUUFEZ2dFQkFMWkE1d2xwSVhDeDRTNE1PaG94TUt5MHVaNGFrSkM4ZW9mdXV2L09ORC9sVWRheXk1QkFObExYCnhISDlCakwvNGczQUk5WWlXQzJuQkpqYUF2dHdhVHBDOTczUldmODhGanV2d1h5Q2YxOGZFb21MbXorK1FiM0IKeitoT3RKZUIvRndVUUJZMWJ4T0RkV3hEUGNieStyUXFPZGJVMEZBUVdmN0NrR2E4NjBGRlhlcHRyRzMrLzFJYQpZVFZ0bm9rZzk5Ky9tRVloVVFFYmFyQWo4aWpmRDVnbnJEYnhERnJWMlg5K1hQOVlJTTlNOVFoY1Z2cGx4WkhlCndFQnYxcHFzSGRHd1hvVnNMRXZjeE9OK1pkZDVsQUlreVlsR3FEenEvc3o3WEJYbWQrSWhYSGdBNWdOZUx4MWQKbWl0UTFlYzVWZjR2SDRJYWMvZjB5RmFTSWhqK2VDdz0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
    client-key-data: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb2dJQkFBS0NBUUVBMWhyVFJ3L3BzaFUyNHRUVkplNk9BT0FKSC9jYjVHR1lLcEVWMmtvUzR0d0QxclFDCk9NTjJhdWgxV3MrTk8wNXBldnBCcTY0NzRFbi9JckxDUkpxTU9NdGJmWEtWMXR1cncwdFg5QXd3bG9mcVo0RlkKTDNpSlNWRDVNT0hJdmtaR2NXemk4eHhIZzF3OWVTNktvZXBGeHR2VklCMFNlOGVRaUFDeGNRdi9iVk5pWEY0NQozS3pQUE1mNG9YKzUwbWgvTWhrUnd0ZCt1S3pibUtoTm84T0FGNE9pRkRwU215SnRjZVVwQUdCSkk0ekw0aEVSCk5SZUFPNVgvaVFkV0U1WmI3aVh6QVV3TVIrKy9wZG00MFZ0NE5mdkwrNjJZdGt2aVZ0NzJ5d2dOSWZhamdiV1IKOVZoR0dhbHJiNmZjTEdrSTB6STlMd2RlWW1Ud3Jod3lWNjI4eXdJREFRQUJBb0lCQUVsSHp5NlFYTDFPRTRZWgpzSWFXR3RaajE5dXYrVVQydmwwN1lVNWdjZ3hobjVLNTg5UzMyZTBIZVR0R3RLRXEybUYwREV3VmkzcmQrTXhJCmdNTjRYaXdHTmw4K2U5aVpRVFhMc09QZjFEV0JlWkpKckFRN2JrbkF4RG1kM2RaNk9Sd1RWNjQ3N0tJaVRRd1EKQ1BVWU9SK3lHRVV3amlYOWpSTnZvVXYrL2tMTldpYWI4NitvNWI5TVZVUVFLVWNOaGI0MTR0ZUVKdjNQaHRUTApKMzdrVllxUWRoUG90VXpPV3BBUkV0dUdmR0Q3ODdGRDFzc3loWmN3dCtENDJiQ0lZbEFhR2FrU3UwY2ZwazBMCmwwWTVnak10enI0L1djOFVUcWM4Ti9GNDlJM3Y0dVN2TDh6dm9nUHBQbjdvZkZ4VTNFblFCL1dXejNrY2hTZUcKN1AzSmowRUNnWUVBNlRBSVlXakNENEVIcERaRUhXTFQ1S2t4SXNVWEZtZXlwZVRSeEpVVjM1L2dvZkFoZ0VsZApuakR5aXQzalFhdlEzdTdpUlJzUDNwVHo4ekJmNlB3VlA3a2pGb2duZmltQ21PVzZBbEVzQTZGUXd5SlBxUUlQCjdqQmY0VFJtbU4vcDhldThjaTJ2VUw3MC9lN0diU0pYZnhRVnd1SnJCYzlON0JrVEpaOUlTU3NDZ1lFQTZ3emgKYTUvL2RXQithaDhYNVloQjF5eUdGMXFHci9GOS8waGdIK0JzbDRWWDY0T09UV09RdEwrWmRrTlhoRGNsT1RyLwpVY2QwYTZNMHBRSFVKUy9idmFXTVFEM2lHRFFsVFpkMkhZbVJ0QUZPWmQrKzREWkZJZlIvbm1LdVRmc1VkTkZVCmdDNm54YVNIOVh3ay9vdXhPdUJtRHZwVHhYelMyNVpCYnk1OVN1RUNnWUFhdHRuKzd0VnNtVWVhMUd5eFFTVXQKU1FUTUN2QStMTnNXamtSSjFScVFaK3lBSU5aMXIvSDlzWFhYMnR1eUpsRGplVktLd0RMdE9QdEpuUDBmMytFLwpUNkpwYm1lMzJHR1J5cCtvckV2eWpvU0dGbVF4NUg2V3R3d0U3TS9rSzZMTmlFZ3FmSGxkTXNNMlpDaHZVRVBkCjF4czNIR0ZPWlJWME90c3FzRGpEeVFLQmdDb0l1NmRMalp1b0RmTmNiQ0dZSkc0ZWdEeGM3QWs2M3BWd2xBQWYKV2t3ZWhJS1JuRGtxdlE0VjFaUFlweVgxMXFwTmFxdHRSOXlYNnlvM0VZRTA5YzhNYy9CcElLM0RaWWhpdGJUQwpqVlByaCtHZ0NicCsrZzRBYzNJWG82USthb0laalVyL0RQSERZcXo3N29HMjZaTGwzbHAvV2N2UGJzWG1NUDE1CkN0OEJBb0dBT25iOFF6czRMOFpFRTRKZzBrNnNDTm1RaUZZVXk0WEc1QTEwUjBibjdUdjJJQWlTaVhuL2l2TmIKVFRxYkZrN1dWY0VLZFlkUm8wQy9Tb0JWeDVOaW9HR3Q2aWN4cGZTY3pMaEZjVWdiQ0orZ1o4TnNtYkZqUFd5ZQpxWGlWc2xTSDRPUjRPV05HZ3gzTGN5OVhwSFdFVEZHMkludVo3Qks2Z21MZG9JQVlSL0k9Ci0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg==
`

// 无缓存new105次： 延迟约5秒
func main() {
	eg.SetLimit(105)
	nums := make([]bool, 105)
	// cache := make([]client.Client, 105)
	// cache[num], _ = newClientFromByte([]byte(kubeconfig))

	start := time.Now()
	for i := range nums {
		num := i
		eg.Go(func() error {
			cli, _ := newClientFromByte([]byte(kubeconfig))
			// 模拟工作中
			err := poll(context.TODO(), 50*time.Millisecond, 10*time.Second, time.Minute, func() (done bool) {
				// cache[num].Scheme()
				cli.Scheme()
				return false
			}, func() error {
				return fmt.Errorf("timeout")
			})
			if err != nil {
				fmt.Println(num, "error:", err)
			}
			return nil
		})
	}
	eg.Wait()
	end := time.Now()
	fmt.Println(end.Sub(start))
}

func main1() {
	eg.SetLimit(105)
	nums := make([]bool, 105)

	start := time.Now()
	for i := range nums {
		num := i
		eg.Go(func() error {
			config := kubeconfig
			id := "id"
			// 模拟config更新
			if n := rand.Intn(1000); n < 10 {
				config += fmt.Sprintf("# %d", n)
			}
			// 模拟不同缓存
			if n := rand.Intn(10); n < 10 {
				id += fmt.Sprint(n)
			}
			isFedControlClusterClientChanged(id, config)
			cli, err := getFedControlClusterClientFromCacheOrSecret(id, config)
			if err != nil {
				fmt.Println(num, "error:", err)
				return nil
			}
			defer deleteFedControlClusterClientIfNotInUse(id)
			// 模拟工作中
			err = poll(context.TODO(), 50*time.Millisecond, 10*time.Second, time.Minute, func() (done bool) {
				cli.Scheme()
				return false
			}, func() error {
				return fmt.Errorf("timeout")
			})
			if err != nil {
				fmt.Println(num, "error:", err)
			}
			return nil
		})
	}
	eg.Wait()
	end := time.Now()
	fmt.Println(end.Sub(start))
}

func getFedControlClusterClientFromCacheOrSecret(id, kubeconfig string) (client.Client, error) {
	lock.Lock()
	defer lock.Unlock()
	if fccClient, ok := fedControlClusterClients[id]; ok {
		fccClient.count++
		if kubeconfig != fccClient.kubeconfig {
			cli, err := newClientFromByte([]byte(kubeconfig))
			if err != nil {
				return nil, err
			}
			fccClient.cli = cli
		}
		fedControlClusterClients[id] = fccClient
		return fccClient.cli, nil
	}
	cli, err := newClientFromByte([]byte(kubeconfig))
	if err != nil {
		return nil, err
	}
	fedControlClusterClients[id] = &fedControlClusterClient{
		count:      1,
		cli:        cli,
		kubeconfig: kubeconfig,
	}
	return cli, nil
}

func deleteFedControlClusterClientIfNotInUse(id string) {
	lock.Lock()
	defer lock.Unlock()
	if fccClient, ok := fedControlClusterClients[id]; ok {
		fccClient.count--
		if fccClient.count == 0 {
			delete(fedControlClusterClients, id)
		} else {
			fedControlClusterClients[id] = fccClient
		}
	}
}

func isFedControlClusterClientChanged(id, kubeconfig string) bool {
	lock.Lock()
	defer lock.Unlock()
	if fccClient, ok := fedControlClusterClients[id]; ok {
		return fccClient.kubeconfig != kubeconfig
	}
	return false
}

func loadKubeConfigFromByte(rawConfig []byte) (*rest.Config, error) {
	apiconfig, err := clientcmd.Load(rawConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed load kubernetes apiconfig")
	}

	clientConfig := clientcmd.NewDefaultClientConfig(*apiconfig, &clientcmd.ConfigOverrides{})
	kubeCfg, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed build client config")
	}

	return kubeCfg, nil
}

func newClient(kubeCfg *rest.Config) (client.Client, error) {
	mapper, err := apiutil.NewDynamicRESTMapper(kubeCfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed new mapper")
	}

	cli, err := client.New(kubeCfg, client.Options{
		Mapper: mapper,
		Scheme: scheme,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed new client")
	}
	return cli, nil
}

func newClientFromByte(kubeconfig []byte) (client.Client, error) {
	cfg, err := loadKubeConfigFromByte(kubeconfig)
	if err != nil {
		return nil, err
	}
	return newClient(cfg)
}

func poll(ctx context.Context, initInterval, maxInterval, timeout time.Duration, condition func() (done bool), timeoutFunc func() error) error {
	loopTimeout := time.After(timeout)
	interval := initInterval
	for {
		select {
		case <-time.After(interval):
			if condition() {
				return nil
			}
			break
		case <-loopTimeout:
			return timeoutFunc()
		case <-ctx.Done():
			return ctx.Err()
		}
		interval *= 2
		if interval > maxInterval {
			interval = maxInterval
		}
	}
}
