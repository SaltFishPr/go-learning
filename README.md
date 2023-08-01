# learning

There is no end to go for learning.

## 协议

- Header: 12 bytes (96 bits)

| Magic Number | Protocol Version | Compress Type | Codec Type | Reserved |
| :----------: | :--------------: | :-----------: | :--------: | :------: |
| 8 bits (42)  |      8 bits      |    4 bits     |   4 bits   | 72 bits  |

- Total Size: 4 bytes

- ServiceMethod

|  Size   | ServiceMethod |
| :-----: | :-----------: |
| 4 bytes |    n bytes    |

- Metadata

|  Size   | Metadata |
| :-----: | :------: |
| 4 bytes | n bytes  |

- Payload

|  Size   | Payload |
| :-----: | :-----: |
| 4 bytes | n bytes |
