Location History [![CircleCI](https://circleci.com/gh/ahmedkamals/location_history.svg?style=svg)](https://circleci.com/gh/ahmedkamals/location_history "Build Status")
================

[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](LICENSE.md "License")
[![release](https://img.shields.io/github/v/release/ahmedkamals/location_history.svg)](https://github.com/ahmedkamals/location_history/releases/latest "Release")
[![codecov](https://codecov.io/gh/ahmedkamals/location_history/branch/main/graph/badge.svg?token=XPINFB5JYV)](https://codecov.io/gh/ahmedkamals/location_history "Code Coverage")
[![GolangCI](https://golangci.com/badges/github.com/ahmedkamals/location_history.svg?style=flat-square)](https://golangci.com/r/github.com/ahmedkamals/location_history "Code Coverage")
[![Go Report Card](https://goreportcard.com/badge/github.com/ahmedkamals/location_history)](https://goreportcard.com/report/github.com/ahmedkamals/location_history "Go Report Card")
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/65feb277726f4a10895f028d460f9196)](https://www.codacy.com/manual/ahmedkamals/location_history?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=ahmedkamals/location_history&amp;utm_campaign=Badge_Grade "Code Quality")
[![GoDoc](https://godoc.org/github.com/ahmedkamals/location_history?status.svg)](https://godoc.org/github.com/ahmedkamals/location_history "Documentation")
[![DepShield Badge](https://depshield.sonatype.org/badges/ahmedkamals/location_history/depshield.svg)](https://depshield.github.io "DepShield")
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fahmedkamals%2Flocation_history.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fahmedkamals%2Flocation_history?ref=badge_shield "Dependencies")

```bash
  _                     _   _               _    _ _     _
 | |                   | | (_)             | |  | (_)   | |
 | |     ___   ___ __ _| |_ _  ___  _ __   | |__| |_ ___| |_ ___  _ __ _   _
 | |    / _ \ / __/ _` | __| |/ _ \| '_ \  |  __  | / __| __/ _ \| '__| | | |
 | |___| (_) | (_| (_| | |_| | (_) | | | | | |  | | \__ \ || (_) | |  | |_| |
 |______\___/ \___\__,_|\__|_|\___/|_| |_| |_|  |_|_|___/\__\___/|_|   \__, |
                                                                        __/ |
                                                                       |___/
```

in-memory location history server.

Table of Contents
-----------------

* [🏎️ Getting Started](#-getting-started)

    * [Prerequisites](#prerequisites)
    * [Installation](#installation)
    * [Examples](#examples)

* [💬 Discussion](#-discussion)

* [🔥 Todo](#-todo)

* [👨‍💻 Credits](#-credits)

* [🆓 License](#-license)

🏎️ Getting Started
------------------

### Prerequisites

* [Golang 1.15 or later][1].

### Installation

```bash
go get -u github.com/ahmedkamals/location_history
cp .env.sample .env
```

### Examples

```bash
make run
```

### ⚓ Git Hooks

In order to set up tests running on each commit do the following steps:

```bash
git config --local core.hooksPath .githooks
```

💬 Discussion
-------------

So the basic idea is using the concept of sliding or rotating window, and let's assume that we have a TTL of 60 seconds.

For every order we would have a slice of 60 placeholders which every value of the type slice would represent a second in
the last minutes (but this is not efficient in terms of memory, we will find a better solution below).

Then we would have a ticker that ticks every second and increases a counter atomically, which you can guess that it
would be used an index for the previous slice, and it would keep rotating over it, and once it visits a position it
overrides whatever data inside.

There is a downside for this solution we rely on subsequent requests to delete the old data, so that means if we stopped
getting requests, we might have a heritage of a memory that is full of ancient legacy locations.

One final not on this, retrieving data in reverse chronological order would be trickier, due to the circular rotation of
the index.

```text
                   Index
 <--- Older data    |
#-------------------x-------------------#
                    <--- Oldest data --
```

So an optimization we can do the following:

* Using a map instead of the slice, so we don't have to pre-allocated the memory and only insert data when necessary.
* Simulating the garbage collector, by applying the same mechanism of mark and sweep technique, so we only turn it on,
  when we detect a cool down in the requests, we sweep the "gold" data.

I did not finish the implementation, because I ran out of time.

<details>
<summary>🔥 Todo:</summary>
   <ul>
       <li>Unit tests.</li>
       <li>Benchmarks.</li>
       <li>Performance optimization.</li>
       <li>Logging command execution.</li>
       <li>Refactoring.</li>
   </ul>
</details>

👨‍💻 Credits
----------

* [ahmedkamals][2]

🆓 LICENSE
----------

Location History is released under MIT license, please refer to
the [`LICENSE.md`](https://github.com/ahmedkamals/location_history/blob/main/LICENSE.md "License") file.

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fahmedkamals%2Flocation_history.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fahmedkamals%2Flocation_history?ref=badge_large "Dependencies")

Happy Coding 🙂

[![Analytics](http://www.google-analytics.com/__utm.gif?utmwv=4&utmn=869876874&utmac=UA-136526477-1&utmcs=ISO-8859-1&utmhn=github.com&utmdt=location_history&utmcn=1&utmr=0&utmp=/ahmedkamals/location_history?utm_source=www.github.com&utm_campaign=location_history&utm_term=location_history&utm_content=location_history&utm_medium=repository&utmac=UA-136526477-1)]()

[1]: https://golang.org/dl/ "Download Golang"

[2]: https://github.com/ahmedkamals "Author"
