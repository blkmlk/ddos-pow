## Word of Wisdom Implementation

### Description

An implementation of a TCP Server that sends a random quote from the Word of Wisdom after receiving a successfully solved challenge

### PoW

I decided to choose the Scrypt algorithm because it's supposed to be GPU-resistant and used for PoW

The leading zeroes check is chosen since I found it commonly used in many publications related to preventing DDoS attacks by using PoW

### How to run?
```shell
make run
```