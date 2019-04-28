# Seek

[![Build Status](https://travis-ci.com/muguangyi/seek.svg?branch=master)](https://travis-ci.com/muguangyi/seek) [![GoDoc](https://godoc.org/github.com/muguangyi/seek?status.svg)](https://godoc.org/github.com/muguangyi/seek) [![Go Report Card](https://goreportcard.com/badge/github.com/muguangyi/seek)](https://goreportcard.com/report/github.com/muguangyi/seek)

**Seek** 是`golang`实现的一套轻量级服务器开发框架，以**信号依赖**为规则建立容器互联，允许用户灵活定制自己的服务器架构，并能快速建立易于扩展的服务器开发解决方案。

## Seek做什么

**Seek**是一套服务器开发架构，规范了服务器开发模式，从而标准化服务间启动，互联和交互流程，使开发者摆脱网络链接，协议等底层逻辑，更关注于模块内逻辑的设计实现，从而开发出更为内聚，松散，通用的功能模块，大大提高复用性。

## Seek不做什么

**Seek**框架中找不到什么`gateway`，`lobby`，`login`等等常见的游戏服务器实现，甚至连`log`模块都没有。这些特定的功能模块**Seek**都不会提供，而是希望开发者基于**Seek**的框架上扩展实现。

## 框架

    +----------------------------+            +=======+  register  +--------------+
    | union                      |  register  |       |<<<<<<>>>>>>| union        |
    |                            |<<<<<<>>>>>>|       |   query    +--------------+
    |                            |   query    |  hub  |
    |                            |            |       |  register  +--------------+
    | +------------------------+ |            |       |<<<<<<>>>>>>| union        |
    | | signal 1               | |            +=======+   query    |              |
    | | signal 2 (book sig N)  | |                                 | +----------+ |
    | +------------------------+ |<------------------------------->| | signal N | |
    |                            |        directly connected       | +----------+ |
    +----------------------------+                                 +--------------+

## 技术点

* 一个信号容器（union）是一个独立的服务器节点，可容纳多个功能信号单元(signal)
* 每一个功能信号单元（signal）运行在一个独立协程中
* 信号与信号间通过管道RPC通信（暂时只提供`同步`方式）
* 不同容器内的信号也可以通过同样的方式通信（建立在**信号依赖**的容器互联）
* **信号基站**：提供信号注册和查询功能，每一个容器都要向至少一个基站注册。

### Signal

所用的功能都是一个独立的`Signal`，并且`Signal`之间的调用无需关心是在同一个容器中，还是不同的容器，是在同一个物理机还是在不同的物理机。

### Union

信号单元容器，是一个独立的服务器节点（可以是一个独立的服务器进程，也允许多个容器在一个服务器进程）。

### Hub

信号容器基站，为信号容器分配端口，管理各种信号的注册和发现功能。