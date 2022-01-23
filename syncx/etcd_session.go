package syncx

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"go.uber.org/multierr"
)

// EtcEtcdSync 代表一个基于Redis的分布式并发原语管理器，底层维护和etcd的连接.
// 通过它可以得到各种分布式的并发原语.
type EtcdSync struct {
	cli     *clientv3.Client
	session *concurrency.Session
}

// NewEtcdSync 返回一个新的EtcdSync对象.
func NewEtcdSync(config clientv3.Config, opts ...concurrency.SessionOption) (*EtcdSync, error) {
	cli, err := clientv3.New(config)
	if err != nil {
		return nil, err
	}

	session, err := concurrency.NewSession(cli)
	if err != nil {
		return nil, err
	}

	return &EtcdSync{cli: cli, session: session}, nil
}

// Close 关闭和底层etcd的连接.
func (d *EtcdSync) Close() error {
	return multierr.Combine(
		d.session.Close(),
		d.cli.Close(),
	)
}
