# netio

Netio is a protocol stack based on golang. It virtualizes the network, so that programs only interact with channels mapped to that network. Since we assume all other connecting nodes run golang and netio channels can map to each other on the other side.
To build any protocol the users have to define message types and a system of channels .

See also 
* the former netchan package
https://groups.google.com/g/golang-nuts/c/Er3TetntSmg/m/1AjSzs_3pzwJ
* https://libp2p.io/
