# tcp_udp_forward

A simple TCP UDP forwarding program.

I noticed that most examples don't return the UDP response from the remote address back to the source address, so I wrote this version.

I hope you don't forget this point during interviews like I did.

---
For example, if I want to forward TCP data.

I can forward local port `8888` to the remote address `192.168.1.1:9999` and local port `7777` to the remote address `10.0.12.1:6666`. 

To start the program, use `./main -p 8888:192.168.1.1:9999/tcp -p 7777:10.0.12.1:6666/tcp`.

---

Alternatively, for forwarding UDP data.

I can forward local port `2024` to the remote address `192.168.1.1:2025`. 

To start the program, use `./main -p 2024:192.168.1.1:2025/udp`.

Of course, you can simultaneously forward TCP and UDP. I believe you already know how to use it.
