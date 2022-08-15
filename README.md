## Crypto-markets arbitrage monitor

Concurrently monitors ByBit, Binance, Garantex buy prices and Garantex sell price. If price difference is bigger than some value, sends notification to discord channel.

### Usage
```
    go build .
    ./monitor
```

### ToDo

 - [x] Fetch ByBit price and order
 - [x] ByBit tests
 - [x] Pretty message in discord with buy link, name, quantity, price for every marketplace 
