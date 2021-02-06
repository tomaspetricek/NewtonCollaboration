# Leader Election
- Follows tutorial: https://medium.com/@felipedutratine/leader-election-in-go-with-etcd-2ca8f3876d79.
- Implements leader election pattern using etcd.
- Have to have etcd running using cmd: etcd.
  
## Working Example
- Shows scenario, when program is run twice at the same time. 
  The first one gets to do its work, while the other has to wait. 
  The program reports back its state into a console.