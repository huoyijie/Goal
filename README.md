# Goal
[Goal](https://github.com/huoyijie/Goal) is a lightweight web framework like Django written in Golang.

[GoalGenerator](https://github.com/huoyijie/GoalGenerator) can auto-generate golang models from yaml file.

[GoalUI](https://github.com/huoyijie/GoalUI) is front of Goal.

## Install Goal

* Go 1.20+

* Repository

```bash
$ git clone git@github.com:huoyijie/Goal.git
```

* Tidy

```bash
$ cd Goal
$ go mod tidy
```

## Run example

* Create Superuser

```bash
$ cd Goal
$ go run ./goalcli -email huoyijie@huoyijie.cn -username huoyijie
```

* Default Database (sqlite)

```bash
$ ls -l ~/.goal/goal.db
```

* Run example `goal-web`

```bash
$ cd Goal/examples
$ go run goal-web.go
```

## Install GoalGenerator

* Repository

```bash
$ go install github.com/huoyijie/GoalGenerator/goalgen@v0.0.29
```

* Generate models from yaml files

```bash
$ cd Goal
$ goalgen -d yamls
```