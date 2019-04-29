# Seek

[![Build Status](https://travis-ci.com/muguangyi/seek.svg?branch=master)](https://travis-ci.com/muguangyi/seek) [![GoDoc](https://godoc.org/github.com/muguangyi/seek?status.svg)](https://godoc.org/github.com/muguangyi/seek) [![codecov](https://codecov.io/gh/muguangyi/seek/branch/master/graph/badge.svg)](https://codecov.io/gh/muguangyi/seek) [![Go Report Card](https://goreportcard.com/badge/github.com/muguangyi/seek)](https://goreportcard.com/report/github.com/muguangyi/seek)

**Seek** is a lightweight server develop framework for `golang`. **Seek** setup connections between containers based on `signal dependency`. This solution provides much flexibility to user to setup extendable servers quickly.

## WHAT Seek DO

**Seek** is a server framework, and define the dev pattern to standalize server startup, connection and communicate flow. So user DO NOT need to write code for low level logic, like network connection, communicate protocol, etc, but only focus on the design and implementation for internal modules. It could make the modules more cohesive, decomposed, and general to improve the reusability.

## WHAT Seek DO NOT do

There is no `gateway`, `lobby`, or `login` server implementation in **Seek**, even no `log` module. Those featured modules will not be provided by **Seek**, but need user to implement based on **Seek** framework.

## Framework Diagram

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

## Tech Notes

* signaler container (union) is an independent server node, and could containe many signals.
* every signal runs within an independent routine.
* communication between signals base on channel RPC (only support sync mode so far)
* signals in different containers could communicate through the same way (RPC based on `signal dependency`)

## Quick Start

## Document