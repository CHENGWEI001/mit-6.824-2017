package raftkv

import "github.com/sunhay/scratchpad/golang/mit-6.824-2017/src/labrpc"
import "crypto/rand"
import "math/big"

type Clerk struct {
	servers []*labrpc.ClientEnd

	lastKnownLeader int
}

func nrand() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	x := bigx.Int64()
	return x
}

func MakeClerk(servers []*labrpc.ClientEnd) *Clerk {
	ck := Clerk{servers: servers}
	return &ck
}

//
// fetch the current value for a key.
// returns "" if the key does not exist.
// keeps trying forever in the face of all other errors.
//
// you can send an RPC with code like this:
// ok := ck.servers[i].Call("RaftKV.Get", &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
//
func (ck *Clerk) Get(key string) string {
	args, reply := GetArgs{Key: key}, GetReply{}

	index := ck.lastKnownLeader
	for reply.Err != OK {
		ok := ck.servers[index%len(ck.servers)].Call("RaftKV.Get", &args, &reply)
		if ok {
			if reply.WrongLeader { // Try next node in server list
				index++
			} else if reply.Err == ErrNoKey { // Return "" if key does not exist
				return ""
			}
		}
	}
	ck.lastKnownLeader = index % len(ck.servers) // Update latest known working server
	return reply.Value
}

//
// shared by Put and Append.
//
// you can send an RPC with code like this:
// ok := ck.servers[i].Call("RaftKV.PutAppend", &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
//
func (ck *Clerk) PutAppend(key string, value string, op string) {
	args, reply := PutAppendArgs{Key: key, Value: value, Op: op}, PutAppendReply{}

	index := ck.lastKnownLeader
	for reply.Err != OK {
		ok := ck.servers[index%len(ck.servers)].Call("RaftKV.PutAppend", &args, &reply)
		if ok {
			if reply.WrongLeader { // Try next node in server list
				index++
			}
		}
	}
	ck.lastKnownLeader = index % len(ck.servers) // Update latest known working server
}

func (ck *Clerk) Put(key string, value string) {
	ck.PutAppend(key, value, "Put")
}
func (ck *Clerk) Append(key string, value string) {
	ck.PutAppend(key, value, "Append")
}
