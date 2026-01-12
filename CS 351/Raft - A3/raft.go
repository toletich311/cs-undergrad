package raft

//
// This is an outline of the API that raft must expose to
// the service (or tester). See comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   Create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   Start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   Each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester) in the same server.
//

import (
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"./labrpc"
)

// As each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make(). Set
// CommandValid to true to indicate that the ApplyMsg contains a newly
// committed log entry.
type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int
}

// make new log struct because I want to use a term indice
type raftLogEntry struct {
	Term    int
	Command interface{}
}

// A Go object implementing a single Raft peer.
// I'm going to use the term "peer" and "server" interchangeable to describe servers is in this implementation***
type Raft struct {
	mu    sync.Mutex          // Lock to protect shared access to this peer's state **
	peers []*labrpc.ClientEnd // RPC end points of all peers
	me    int                 // This peer's index into peers[]
	dead  int32               // Set by Kill()

	// Your data here (3A, 3B).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.

	// data for persistent state on all servers - from figure 2
	currentTerm int            //latest term server has seen (initialized to 0 on first boot, increases monotonically)
	votedFor    int            //candidateId that received vote in current term (or null if none)
	log         []raftLogEntry //log entries; each entry contains command for state machine, and term when entry was received by leader (first index is 1)

	// data for volatile state on all servers - from figure 2
	commitIndex int //index of highest log entry known to be committed (initialized to 0, increases monotonically)
	lastApplied int //index of highest log entry applied to statemachine (initialized to 0, increases monotonically)

	// data for volatile state on all leaders - from figure 2
	// Reinitialized after election
	nextIndex  []int //for each server, index of the next log entry to send to that server (initialized to leader last log index + 1)
	matchIndex []int //for each server, index of highest log entry known to be replicated on server (initialized to 0, increases monotonically)

	// data for election - I created these because I know I'll need to keep track
	state           string        //whether each peer is a leader, follower, or candidate
	electionTimeout time.Duration //a random election timeout per peer
	lastHeartbeat   time.Time     //a timer counted the time since they recieved the last hearbeat - (will later compare to electiontimeout to determine if the peer will start a new election)

	// adding channel for application message!!!!
	applyCh chan ApplyMsg
}

// Return currentTerm and whether this server
// believes it is the leader.
// use this function to getstate() during consensus to avoid race conditions and deadlock!!
func (rf *Raft) GetState() (int, bool) {
	//first lock and defer unlock to prevent race conditions
	rf.mu.Lock()
	defer rf.mu.Unlock()                        //i'm using defer in case it takes time to retrieve currentTerm or state fields in raft struct
	return rf.currentTerm, rf.state == "Leader" //basically return specified instance of raft structs current state
}

// Example RequestVote RPC arguments structure.
// Field names must start with capital letters!
type RequestVoteArgs struct {
	// Your data here (3A, 3B).
	// add more data ---
	// -----term, candidateid, lastlogidx, lastlogterm
	// these arguments come from raft paper figure 2
	Term         int //candidate’s term
	CandidateId  int //candidate requesting vote
	LastLogIndex int //index of candidate’s last log entry (§5.4)
	LastLogTerm  int //term of candidate’s last log entry (§5.4)
}

// Example RequestVote RPC reply structure.
// Field names must start with capital letters!
type RequestVoteReply struct {
	// Your data here (3A).
	// add more data ---
	// -------term, voteGranted
	// these results come from raft paper figure 2
	Term        int  //currentTerm, for candidate to update itself
	VoteGranted bool //true means candidate received vote
}

// Example RequestVote RPC handler.
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	// Your code here (3A, 3B).
	// lock rf.mu.Lock() then defer unlock
	// -- this is to prevent race conditions on rf fields
	rf.mu.Lock()
	defer rf.mu.Unlock()

	//reject outdated candidate terms
	// Reply false if term < currentTerm (§5.1) (raft figure 2)
	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		reply.VoteGranted = false
		return
	}

	//in the other case we ensure that everyone matches the highest candidate term
	if args.Term > rf.currentTerm {
		rf.currentTerm = args.Term
		rf.votedFor = -1
	}

	//if first vote, cast vote , else don't vote
	//check if log is up to date
	//--If votedFor is null or candidateId, and candidate’s log is at least as up-to-date as receiver’s log, grant vote (§5.2, §5.4) (raft figure 2)
	if (rf.votedFor == -1 || rf.votedFor == args.CandidateId) && rf.isLogUpToDate(args.LastLogIndex, args.LastLogTerm) {
		rf.votedFor = args.CandidateId
		reply.VoteGranted = true
	} else {
		reply.VoteGranted = false
	}

	//set reply term to current term
	reply.Term = rf.currentTerm
}

// Example code to send a RequestVote RPC to a server.
// Server is the index of the target server in rf.peers[].
// Expects RPC arguments in args. Fills in *reply with RPC reply,
// so caller should pass &reply.
//
// The types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
//
// Call() is guaranteed to return (perhaps after a delay) *except* if the
// handler function on the server side does not return.  Thus there
// is no need to implement your own timeouts around Call().
//
// Look at the comments in ../labrpc/labrpc.go for more details.
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	return ok
}

// The service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. If this
// server isn't the leader, returns false. Otherwise start the
// agreement and return immediately. There is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election. Even if the Raft instance has been killed,
// this function should return gracefully.
//
// The first return value is the index that the command will appear at
// if it's ever committed. The second return value is the current
// term. The third return value is true if this server believes it is
// the leader.
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	log.Println("===> Start() called") //logging to fix deadlock

	/*	Called by the service when it wants to start agreement on a new command.
		Only the leader should accept it, append it to its log, and start replication.
		Need to: */
	//- Check if rf is the leader
	term, isLeader := rf.GetState() //use getstate to prevent race condition and deadlock

	//if not leader return
	if !isLeader {
		return -1, -1, false
	}

	entry := raftLogEntry{
		Term:    term,
		Command: command,
	}

	rf.mu.Lock() //preventing race condition since we acess command and term

	//replicate the new entry
	rf.log = append(rf.log, entry)

	//- Return (index, term, isLeader)
	index := len(rf.log) - 1                                                                   // log is 0-indexed, but Raft starts log indices at 1
	log.Printf("===> Start(): Appended command %v at index %d, term %d", command, index, term) //logging to fix overwriting log issue
	rf.mu.Unlock()
	return index, term, true

}

// The tester doesn't halt goroutines created by Raft after each test,
// but it does call the Kill() method. Your code can use killed() to
// check whether Kill() has been called. The use of atomic avoids the
// need for a lock.
//
// The issue is that long-running goroutines use memory and may chew
// up CPU time, perhaps causing later tests to fail and generating
// confusing debug output. Any goroutine with a long-running loop
// should call killed() to check whether it should stop.
func (rf *Raft) Kill() {
	atomic.StoreInt32(&rf.dead, 1)
	// Your code here, if desired.
}

func (rf *Raft) killed() bool {
	z := atomic.LoadInt32(&rf.dead)
	return z == 1
}

// The ticker go routine starts a new election if this peer hasn't received
// heartsbeats recently.
func (rf *Raft) ticker() {
	for !rf.killed() { //loop runs while the server is up

		//sleep for random duration
		time.Sleep(10 * time.Millisecond)

		//safe comparison before starting go routines
		//- here we check if the current heartbeat is taking longer than the election timeout
		rf.mu.Lock()
		timeout := time.Since(rf.lastHeartbeat) > rf.electionTimeout
		rf.mu.Unlock()

		//if not leader and ,equal electiontimout -- start election (with helper func)
		// - here we see that if a server is not the leader and timed out, we start a new election
		rf.mu.Lock()
		if rf.state != "Leader" && timeout {
			//adding logging
			log.Printf("[Peer %d] Election timeout triggered after %v ms", rf.me, time.Since(rf.lastHeartbeat).Milliseconds())
			go rf.startElection() //this go routine starts an election
		}
		rf.mu.Unlock()
	}
}

// The service or tester wants to create a Raft server. The ports
// of all the Raft servers (including this one) are in peers[]. This
// server's port is peers[me]. All the servers' peers[] arrays
// have the same order. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
func Make(peers []*labrpc.ClientEnd, me int, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.me = me

	//Step 1 - intialize fields from raft struct that i defined from raft figure 2
	//--currentterm, voted for, log
	rf.mu.Lock()                       //prevent race condition from current term and log access
	rf.currentTerm = 0                 //all servers start at term 0
	rf.votedFor = -1                   // all servers/peers have not voted intially
	rf.log = []raftLogEntry{{Term: 0}} // alll logs are empty (with dummy node of 0 to indicate so)
	rf.mu.Unlock()

	//initalize the beginning state
	rf.mu.Lock()          //added mutex's because state is accessed throughout goroutines so every read/write needs to be in a lock!!!!
	rf.state = "Follower" //all servers start out as followers
	rf.mu.Unlock()

	//--electiontimout, heartbeat
	rf.electionTimeout = time.Duration(300+rand.Intn(300)) * time.Millisecond //every server has a random election timeout
	rf.lastHeartbeat = time.Now()                                             //initilaize the first heartbeat to now, so all servers will start first election at different times!

	// initialze channel to take application message here
	rf.applyCh = applyCh

	// start ticker goroutine to start elections.
	// -- ticker is going to use the time fields above to elect/relect only when neccessary
	go rf.ticker()
	go rf.applyEntries()

	return rf
}

// Create helper functions
// --isLogUpToDate()-- this helper function compares log indices and terms to decide if the candidate is eligible
func (rf *Raft) isLogUpToDate(candidateLastIndex int, candidateLastTerm int) bool {
	//Im creating vars lastIndex and last term to make comparisons easier to read in this function
	lastIndex := len(rf.log) - 1       // we used a dummy node so we have to save the last index as the length -1
	lastTerm := rf.log[lastIndex].Term //here i save the lastTerm in a seperate var

	//to determine if a log is up to date we check if the candidate's last term is greater  or equal to  the followers'
	if candidateLastTerm > lastTerm {
		return true
	} else if candidateLastTerm == lastTerm {
		return candidateLastIndex >= lastIndex //should also always return true
	} else {
		return false //otherwise we say this candidate is behind and cannot get a vote from the server, so we return false
	}
}

// --startElection()
// pseudocode based on figure 2 of raft
func (rf *Raft) startElection() {
	//reset the last heartbeat -- this helped avoid servers continously starting elections
	// added during debugging
	rf.lastHeartbeat = time.Now()

	//lock to prevent datarace - rf state field is accessed in mutliple goroutines and requires mutex!
	rf.mu.Lock()

	//initialize rf fields to officially select this server as a candidate
	rf.state = "Candidate"        // i am a candidate
	rf.currentTerm++              // i increase my term
	rf.votedFor = rf.me           // i vote for myself
	votesReceived := 1            // i vote for myself and count it
	termAtStart := rf.currentTerm // here i create a variable to take note of the term i started with for this election

	//added logging during debugging to determine why election timeouts kept occuring
	log.Printf("[Peer %d] Starting election for term %d", rf.me, termAtStart)

	rf.mu.Unlock() // Unlock before sending RPCs

	//for i in range rf.peers
	//- here we loop through all the peers of this server, we are going to request that each one votes for this candidate
	for i := range rf.peers {
		//--if i != rf.me
		// if this peer server is NOT myself (i already voted for myself)
		if i != rf.me {
			//-- go func(peer int)
			// here i send out a goroutine to request a vote, this way the votes can run concurrently
			// however during debugging i needed to add mutex's above because here I access terms and candidateid and state from multiple
			//  concurrent go routines
			go func(peer int) {
				//here i use the requestVoteArgs to take not of all this arguements, this makes comparisons later easier
				args := RequestVoteArgs{
					Term:         termAtStart, // ← store a local copy from rf.currentTerm
					CandidateId:  rf.me,
					LastLogIndex: len(rf.log) - 1,
					LastLogTerm:  rf.log[len(rf.log)-1].Term,
				}

				//send request to get reply
				reply := RequestVoteReply{}
				//logging debug statement - let's me know that my candidate send the request to get a vote out
				log.Printf("[Peer %d] Sending RequestVote to Peer %d for term %d", rf.me, peer, termAtStart)

				//-- if rf.sendRequestVote
				// here we check if we got a response from the other server, it will tell us which server it is, their info, and whether they voted for us or not
				if rf.sendRequestVote(peer, &args, &reply) {

					//added locks during debugging because of concurrent goroutines accessing fields in rf struct.
					rf.mu.Lock()
					defer rf.mu.Unlock()

					//added log statement during debugging to let me know that this candidate got a response/vote back from the specified follower server
					if rf.state == "Candidate" && rf.currentTerm == args.Term && reply.VoteGranted {
						log.Printf("[Peer %d] Vote granted from Peer %d", rf.me, peer)
					}

					//check for term change safely
					// If term changed, revert to follower
					// this is an instance where the candidate that want to be a leader is out of date, in this instance this candidate should not be a leader
					if reply.Term > rf.currentTerm {
						rf.currentTerm = reply.Term
						rf.votedFor = -1
						rf.state = "Follower"
						return
					}

					//---------if votesreceived >len rf.peers/2
					//--------become leader
					// here we check if first if this server hasn't been reverted to a follower, if its up to date, and if it received vote
					// in this instance, we increment votes recieved
					if rf.state == "Candidate" && rf.currentTerm == args.Term && reply.VoteGranted {
						votesReceived++
						//each time a vote is recieved we check if the candidate has recieved a majority of votes
						if votesReceived > len(rf.peers)/2 {
							//added log statements to debug my issue of elections timing out
							log.Printf("[Peer %d] Became Leader for term %d", rf.me, rf.currentTerm)
							//in this insance the candidate becomes the leader
							rf.becomeLeader()
						}
					}
				}

				//log getting reply - this was for debugging because my elections timed out every time
				log.Printf("[Peer %d] Received vote reply from Peer %d: VoteGranted=%v, Term=%d", rf.me, peer, reply.VoteGranted, reply.Term)

			}(i)
		}
	}
}

// --becomeLeader()
func (rf *Raft) becomeLeader() {
	//when a candidate becomes leader, we change the state and send heartbeats
	// had to delete mutex's - the mutexs in startElection already prevent race conditions for rf.state write here
	// the mutex's I had here before led to a deadlock, which prevented hearbeats from being sent-- this was the root of the issue of elections timing out every single time
	rf.state = "Leader"

	// From Raft Figure 2 (volatile state on leaders):
	// - nextIndex[i] = index of the next log entry to send to peer i (initialize to last log index + 1)
	// - matchIndex[i] = index of highest log entry known to be replicated on peer i (initialize to 0)
	// logic for 3b -- initialize lastlog, next, match idx

	LastLogIndex := len(rf.log) - 1 //log starts at index 0
	rf.nextIndex = make([]int, len(rf.peers))
	rf.matchIndex = make([]int, len(rf.peers))

	for i := range rf.peers {
		rf.nextIndex[i] = LastLogIndex + 1 // send next entry
		rf.matchIndex[i] = 0               // none to be replicated
	}

	//send out heartbeats as leader
	go rf.sendHeartbeats()
}

// add structs for appendEntries
// took arguments from raft paper figure 2
type AppendEntriesArgs struct {
	Term     int
	LeaderId int

	//add more fields based on figure 2 for 3b
	PrevLogIndex int            // index of log entry immediately preceding new ones
	PrevLogTerm  int            // term of prevLogIndex entry
	Entries      []raftLogEntry // log entries to store (empty for heartbeat)
	LeaderCommit int            // leader’s commitIndex

}
type AppendEntriesReply struct {
	Term    int
	Success bool
}

// rpc handler for now
// called whenever heartbeat is sent
func (rf *Raft) sendAppendEntries(server int, args *AppendEntriesArgs, reply *AppendEntriesReply) bool {
	//use channel to make rpc call to gather replies from peer servers
	ok := rf.peers[server].Call("Raft.AppendEntries", args, reply)
	//adding logging because elections always timed out - maybe heartbeats were never being sent
	log.Printf("[Leader %d] AppendEntries RPC to Peer %d returned %v", rf.me, server, ok)
	return ok
}

// helper  for rpc handler -- appendEntries
func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	//added logging for recieving heartbeat because elections always timed out - maybe heartbeats were never recieved
	log.Printf("[Peer %d] Received heartbeat from Leader %d (term %d)", rf.me, args.LeaderId, args.Term)

	//updating heartbeat timer because heartbeat was recieved -- this was another root cause of my elections always timing out!!!!
	rf.lastHeartbeat = time.Now()

	//response to actual heartbeat
	if args.Term > rf.currentTerm { // case 1 - leader may have crashed or got out of date and booted , now a follower
		log.Printf("[Peer %d] Updating term from %d to %d, reverting to Follower", rf.me, rf.currentTerm, args.Term)
		rf.currentTerm = args.Term
		rf.votedFor = -1
		rf.state = "Follower"
	} else if args.Term < rf.currentTerm { //case 2 -  got a heartbeat, successfully return
		reply.Term = rf.currentTerm
		reply.Success = false
		return
	} else { //another edge case where term matches, but another server was elected leader, must revert to follower
		if rf.state != "Follower" {
			log.Printf("[Peer %d] Term matches but was %s. Reverting to Follower", rf.me, rf.state)
			rf.state = "Follower"
		}
	}

	//from figure 2 -- need to add logic for if term<currentTerm to achieve consensus
	/* Current plan, Update to: */

	//1. Reply false if term < currentTerm
	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		reply.Success = false
		return
	}
	// otherwise revert to follower r--> if args.Term > rf.currentTerm
	if args.Term > rf.currentTerm {
		rf.currentTerm = args.Term
		rf.votedFor = -1
		rf.state = "Follower"
	}

	//2. Reply false if log doesn’t contain entry at prevLogIndex whose term matches prevLogTerm
	if args.PrevLogIndex >= len(rf.log) || rf.log[args.PrevLogIndex].Term != args.PrevLogTerm {
		reply.Term = rf.currentTerm
		reply.Success = false
		return
	}

	//3. If an existing entry conflicts with a new one (same index but different term), delete it and all that follow it
	for i, entry := range args.Entries {
		index := args.PrevLogIndex + 1 + i
		if index < len(rf.log) {
			log.Printf("[Follower %d] Log conflict at index %d: existing term=%d, new term=%d",
				rf.me, index, rf.log[index].Term, entry.Term)
			if rf.log[index].Term != entry.Term {
				if index <= rf.commitIndex {
					log.Fatalf("[Follower %d] Attempting to overwrite committed entry at index %d", rf.me, index)
				}
				// append new entries
				rf.log = rf.log[:index]
				rf.log = append(rf.log, args.Entries[i:]...)
				break
			}
		} else {
			//4. Append any new entries not already in the log
			rf.log = append(rf.log, args.Entries[i:]...)
			break
		}
	}
	//5. If leaderCommit > commitIndex, set commitIndex = min(leaderCommit, index of last new entry)
	lastNewIndex := args.PrevLogIndex + len(args.Entries)
	if args.LeaderCommit > rf.commitIndex {
		if lastNewIndex < args.LeaderCommit {
			rf.commitIndex = lastNewIndex
		} else {
			rf.commitIndex = args.LeaderCommit
		}
	}
	log.Printf("[Follower %d] AppendEntries: PrevLogIndex=%d PrevLogTerm=%d, MyLogLen=%d", rf.me, args.PrevLogIndex, args.PrevLogTerm, len(rf.log))

	//6. Send back reply.Term and Success */
	rf.lastHeartbeat = time.Now()
	reply.Term = rf.currentTerm
	reply.Success = true
}

// --sendHeartbeats() this is a helper function called in goroutine by leader to get heartbeats from all its followers
func (rf *Raft) sendHeartbeats() {

	//while leader - this loop keeps sending heartbeats every 100 millisenconds
	for {

		//added mutex to protect rf. state , only send heartbeat if you are a leader, double check here
		rf.mu.Lock()
		if rf.state != "Leader" {
			rf.mu.Unlock() //unlock before returning
			return
		}

		//now we know we are a leader, so we store the current term to use for logging later
		//3B --> intialize volatile state fields
		commitIndex := rf.commitIndex
		currentTerm := rf.currentTerm
		rf.mu.Unlock() //unlock, no longer accessing rf.state

		//loop through all the leader's followers
		for i := range rf.peers {
			//--if i != rf.me
			// if this peer isn't myself (i only send heartbeats to my followers as leader), then i send heartbeat
			if i != rf.me {
				//---- use go function and appendentries args, send, and reply to send heartbeats
				go func(peer int) {
					//Maybe lock here
					rf.mu.Lock() //please don't lead to a deadlock

					//3B logic -- prepare entries
					//	- For each peer, send AppendEntries with actual log entries starting at nextIndex[i]
					prevLogIndex := rf.nextIndex[peer] - 1
					prevLogTerm := rf.log[prevLogIndex].Term

					//replicate logs
					entries := make([]raftLogEntry, len(rf.log[rf.nextIndex[peer]:]))
					copy(entries, rf.log[rf.nextIndex[peer]:])

					//use args from raft paper fig 2
					args := AppendEntriesArgs{
						Term:         currentTerm,
						LeaderId:     rf.me,
						PrevLogIndex: prevLogIndex,
						PrevLogTerm:  prevLogTerm,
						Entries:      entries,
						LeaderCommit: commitIndex,
					}

					//logging to fix overwriting commited values issue
					log.Printf("[Leader %d] Sending AppendEntries to Peer %d: PrevLogIndex=%d PrevLogTerm=%d Entries=%v Commit=%d",
						rf.me, peer, args.PrevLogIndex, args.PrevLogTerm, args.Entries, args.LeaderCommit)

					//now we are done accessessing rf struct fields and therefore can safely unlock
					rf.mu.Unlock()

					//store reply
					reply := AppendEntriesReply{}

					//logging to see if got reply
					log.Printf("[Leader %d] Received AppendEntries reply from %d: Success=%v, Term=%d", rf.me, peer, reply.Success, reply.Term)

					//now send heartbeat
					ok := rf.sendAppendEntries(peer, &args, &reply)

					//add logging for sending a heartbeat - maybe the heartbeat send was never reached
					log.Printf("[Leader %d] Sending heartbeat to Peer %d (term %d)", rf.me, peer, currentTerm)

					//error checking -recieved no entries
					if !ok {
						//maybe add logging here if testbasic fails
						log.Printf("[Leader %d] AppendEntries to %d failed, decrementing nextIndex to %d", rf.me, peer, rf.nextIndex[peer])

						return
					}

					//at this point of recieving successful entries, we are going to access rf. fields again and therefore require a mutex again
					//lock
					rf.mu.Lock()
					defer rf.mu.Unlock() //please don't create a deadlock

					// do I need to step down as leader??? - special case
					if reply.Term > rf.currentTerm {
						rf.currentTerm = reply.Term
						rf.state = "Follower"
						rf.votedFor = -1
						return
					}

					// After receiving successful AppendEntries replies:
					if reply.Success {

						// Update matchIndex and nextIndex
						rf.matchIndex[peer] = args.PrevLogIndex + len(args.Entries)
						rf.nextIndex[peer] = rf.matchIndex[peer] + 1

						// Try to advance commitIndex if possible
						for i := rf.commitIndex + 1; i < len(rf.log); i++ {
							count := 1 // count self

							//- For each log index i > commitIndex:
							// Check if a majority of matchIndex[] ≥ i
							for j := range rf.peers {
								if j != rf.me && rf.matchIndex[j] >= i {
									count++
								}
							}
							// - Check if a majority of matchIndex[] ≥ i AND log[i].Term == currentTerm
							if count > len(rf.peers)/2 && rf.log[i].Term == rf.currentTerm {
								// - Then set commitIndex = i
								rf.commitIndex = i
							}
						}
					} else {
						// Backtrack nextIndex on failure
						if rf.nextIndex[peer] > 1 {
							rf.nextIndex[peer]--

							//send another go routine
							//go rf.sendHeartbeats()

						}
					}
				}(i)
			}
		}

		//sleep for heartbeat time (100 milliseconds) - heartbeat timer
		time.Sleep(100 * time.Millisecond)
	}

}

//probable helper functions that will be needed for consensus

func (rf *Raft) applyEntries() {

	// In a loop:
	for {
		time.Sleep(10 * time.Millisecond) // avoid tight spinning

		rf.mu.Lock()
		// - Check if commitIndex > lastApplied
		for rf.commitIndex > rf.lastApplied {
			// - If yes, increment lastApplied
			rf.lastApplied++
			// - Send ApplyMsg{CommandValid: true, Command: rf.log[lastApplied].Command, CommandIndex: lastApplied} to applyCh
			entry := rf.log[rf.lastApplied]
			msg := ApplyMsg{
				CommandValid: true,
				Command:      entry.Command,
				CommandIndex: rf.lastApplied,
			}
			rf.mu.Unlock() // Unlock before sending to avoid blocking other goroutines

			// - Apply all entries up to commitIndex by sending them over applyCh???

			rf.applyCh <- msg

			rf.mu.Lock() // Re-lock before checking the next entry
			log.Printf("[Peer %d] Applying log index %d: command=%v", rf.me, rf.lastApplied, entry.Command)
		}
		rf.mu.Unlock()
	}
}

/*

getLogTerm(index int) int – safely return log term at index
matchLog(prevLogIndex, prevLogTerm) bool – helper for AppendEntries match check
truncateLogFrom(index int) – delete conflicting entries
*/
