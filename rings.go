////////////////////////////////////////////////////////////////////////////////
// Program: rings.go                                                           /
// Date: 10/31/16                                                              /
// Author: Cody Kankel                                                         /
// Description: rings.go is a Go program developed for B424 Parallel Prog.     /
// class. Rings.go will send a message around a ring of MPI processes 10 times /
// and exit after that. Go does not have native MPI bindings, and so this prog /
// uses the go-mpi bindings found here: https://github.com/JohannWeging/go-mpi /
////////////////////////////////////////////////////////////////////////////////

//NOTE: This program was developed for a class assignment, and was uploaded to Github
// *after* submission and is not intended to divulge any answers or short cuts.
// (The default programming lanugage for the class is C anyways!)

package main

import (
	"fmt"
	"os"
	mpi "github.com/JohannWeging/go-mpi"
)

func main() {
	mpi.Init(&os.Args)
	// Ignoring error messages from func calls
	world_size, _ := mpi.Comm_size(mpi.COMM_WORLD)
	rank, _ := mpi.Comm_rank(mpi.COMM_WORLD)
	//starting out with count as 0. Stopping when its 11 (that means it wrapped around 10 times)
	count := 0

	for count < 11 {
                if rank == 0 {
                        // Increment the counter, as this indicates we wrapped back around to the beginning.
			count++
			mpi.Send( &count, 1, mpi.INT, 1, 0, mpi.COMM_WORLD)
			fmt.Println("Rank 0 just sent count to rank 1")
			if count == 11 {
				//If the count is at 11, we don't want to recv any more messages.
				fmt.Println("Rank 0 is now Finalizing Mpi.")
				mpi.Finalize()
			} else {
				//fmt.Println("Rank 0 is waiting to receive from world -1...")
				mpi.Recv( &count, 1, mpi.INT,( world_size -1), 0, mpi.COMM_WORLD)
				fmt.Println("Rank 0 has received count of ", count, " from world-1!")
			}
		} else {
			if count == 10 {
				// IF the count is at 10, we want to pass the message and finalize
				mpi.Recv( &count, 1, mpi.INT, (rank -1), 0, mpi.COMM_WORLD)
                                mpi.Send( &count, 1, mpi.INT, ((rank + 1) % world_size), 0, mpi.COMM_WORLD)
				fmt.Println("I, rank ", rank, " am finalizing MPI.")
				mpi.Finalize()
			} else if count < 10 {
				mpi.Recv( &count, 1, mpi.INT, (rank -1), 0, mpi.COMM_WORLD)
				fmt.Println("Rank ", rank, " just received a count of ", count)
				mpi.Send( &count, 1, mpi.INT, ((rank + 1)% world_size), 0, mpi.COMM_WORLD)
			}
		}
	}
}

