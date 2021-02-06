replace (
	go.etcd.io/etcd/api/v3 v3.5.0-pre => go.etcd.io/etcd/api/v3 v3.0.0-20210107172604-c632042bb96c
	go.etcd.io/etcd/pkg/v3 v3.5.0-pre => go.etcd.io/etcd/pkg/v3 v3.0.0-20210107172604-c632042bb96c
)

module LeaderSelection

go 1.15

require go.etcd.io/etcd/client/v3 v3.0.0-20210204162551-dae29bb719dd
