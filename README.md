# gofer

[![License](https://img.shields.io/:license-apache%202-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![GoDoc](https://godoc.org/github.com/smallnest/gofer?status.png)](http://godoc.org/github.com/smallnest/gofer)  [![build](https://github.com/smallnest/gofer/actions/workflows/build.yml/badge.svg)](https://github.com/smallnest/gofer/actions/workflows/build.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/smallnest/gofer)](https://goreportcard.com/report/github.com/smallnest/gofer) [![coveralls](https://coveralls.io/repos/smallnest/gofer/badge.svg?branch=master&service=github)](https://coveralls.io/github/smallnest/gofer?branch=master)


企业级中间件框架，Go生态圈的spring框架。

互联网大厂的多年的丰富经验，解决使用Go语言开发企业级产品的痛点，提炼和完善的一套健壮、易用的企业级框架。

目标是类似Java生态圈最知名的 spring框架，成为Go生态圈的应用框架。但是和spring框架完全不同，不会基于依赖注入等模式实现面向对象的go-spring框架，而是采用地道的Go惯用法，在成熟的Go生态群库和自研的在生产中经受考验的总结的库基础上，提供一套常用的企业应用框架.


## 目前提供的功能

- mq: 消息队列读写库
  - mka: 支持多kafka集群的读写，用于容错。
    大家使用kafka的最大的痛点是什么？莫名其妙的kafka不可用，或者某个机房、Region网络故障不可写。这个库就专门解决kafka的容错功能，采用多Kafka集群的方式，总能保证一个可用的Kafka集群.
- syncx: 分布式和扩展的并发原语,包括：
  - Locker: 实现了sync.Locker
  - Mutex: 分布式的锁
  - RWMutex: 分布式的读写锁
  - Barrier: 栅栏
  - DoubleBarrier: 双次栅栏
  - Election: 选主