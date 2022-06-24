package nats

//
// PUB <subject> <size>\r\n
// <message>\r\n
//
// SUB <subject> <sid>\r\n
// SUB <subject> <queue> <sid>
//
// MSG <subject> <sid> <size>\r\n
// <message>\r\n

// https://github.com/jhunters/bigqueue
// https://github.com/pedrogao/trie
