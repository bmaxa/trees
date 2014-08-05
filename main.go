package main

import (
  "fmt"
  "github.com/bmaxa/trees/treap"
  "github.com/bmaxa/trees/aa"
  "github.com/bmaxa/trees/rb"
  "github.com/bmaxa/trees/generic"
  "math/rand"
  "time"
  . "testing"
  "runtime"
  "os"
  "strconv"
)

var N = 1000000
var  rnd1,rnd2,rnd3 []int

func init(){
  if len(os.Args)>1 {
    N,_ = strconv.Atoi(os.Args[1])
  }
  fmt.Println(N)
  rnd1,rnd2,rnd3 = RandomShuffle(N),RandomShuffle(N),array(N)
}

//var  rnd1,rnd2 = array(N),RandomShuffle(N)

func RandomShuffle(n int)[]int{
  r := rand.New(rand.NewSource(time.Now().UnixNano()))
  fmt.Println("shuffle")
  return r.Perm(n)
}

func array(n int)[]int {
  arr := make([]int,n)
  for i:=0;i<n;i++{
    arr[i]=i
  }
  return arr
}

type Key int32
func (t Key) Less(k tree.Key) bool {
  return t<k.(Key)
}

func bench(t tree.ITree,rnd1,rnd2[]int,header string) (inser,find,iter,ins,delet int64){
  fmt.Printf("%s\n--------------------------------\n",header)
  inser = Benchmark(func (b *B) {
  b.N = N
  t.Clear()
  b.ResetTimer()
  var item tree.Item
  for i,k := range rnd1 {
    item.Key = Key(k)
    item.Value = N-k
    t.Insert(item)
    t.Delete(Key(rnd2[i]))
  }
  }).NsPerOp()
  find = Benchmark(func (b *B) {
  b.N = N
  b.ResetTimer()
  for _,k := range rnd1 {
    t.Find(Key(k))
  }
  }).NsPerOp()
  
  iter = Benchmark(func (b *B) {
  b.N = N
  b.ResetTimer()
  for i := t.Begin(); i != t.End(); i=i.Next() {
  }
  for i := t.RBegin(); i != t.End(); i=i.Prev() {
  }
  }).NsPerOp()
  left,right := t.Weight()
  fmt.Printf("height : %v\nweight : (%d,%d)\nvalid : %v\nsize %d\n",t.Height(),left,right,t.Validate(),t.Size())
  t.Clear()
  ins = Benchmark(func (b *B) {
  b.N = N
  t.Clear()
  b.ResetTimer()
  var item tree.Item
  for _,k := range rnd1 {
    item.Key = Key(k)
    item.Value = N-k
    t.Insert(item)
  }
  }).NsPerOp()
  left,right = t.Weight()
  fmt.Printf("height : %v\nweight : (%d,%d)\nvalid : %v\nsize %d\n",t.Height(),left,right,t.Validate(),t.Size())
  delet = Benchmark(func (b *B) {
  b.N = N
  var item tree.Item
  for _,k := range rnd1 {
    item.Key = Key(k)
    item.Value = N-k
    t.Insert(item)
  }
  b.ResetTimer()
  for _,k := range rnd1 {
    t.Delete(Key(k))
  }
  }).NsPerOp()
  return
}

func main() {
  t := rb.New()
  rb_inser,rb_find,rb_iter,rb_ins,rb_delet := bench(t,rnd1,rnd2,"Red Black")
  tt := aa.New()
  aa_inser,aa_find,aa_iter,aa_ins,aa_delet := bench(tt,rnd1,rnd2,"AA")
  ttt := treap.New(nil)
  tr_inser,tr_find,tr_iter,tr_ins,tr_delet := bench(ttt,rnd1,rnd2,"Treap")
  ms := runtime.MemStats{}
  fmt.Printf("Tree\t\tInsert/Erase\tFind\t\tIter\t\tInsert\t\tDelete\n")
  fmt.Printf("RB\t\t%d\t\t%d\t\t%d\t\t%d\t\t%d\n",rb_inser,rb_find,rb_iter,rb_ins,rb_delet)
  fmt.Printf("AA\t\t%d\t\t%d\t\t%d\t\t%d\t\t%d\n",aa_inser,aa_find,aa_iter,aa_ins,aa_delet)
  fmt.Printf("TR\t\t%d\t\t%d\t\t%d\t\t%d\t\t%d\n",tr_inser,tr_find,tr_iter,tr_ins,tr_delet)
  rb_inser,rb_find,rb_iter,rb_ins,rb_delet = bench(t,rnd3,rnd2,"Red Black")
  aa_inser,aa_find,aa_iter,aa_ins,aa_delet = bench(tt,rnd3,rnd2,"AA")
  tr_inser,tr_find,tr_iter,tr_ins,tr_delet = bench(ttt,rnd3,rnd2,"Treap")
  fmt.Println("Ordered Insert")
  fmt.Printf("Tree\t\tInsert/Erase\tFind\t\tIter\t\tInsert\t\tDelete\n")
  fmt.Printf("RB\t\t%d\t\t%d\t\t%d\t\t%d\t\t%d\n",rb_inser,rb_find,rb_iter,rb_ins,rb_delet)
  fmt.Printf("AA\t\t%d\t\t%d\t\t%d\t\t%d\t\t%d\n",aa_inser,aa_find,aa_iter,aa_ins,aa_delet)
  fmt.Printf("TR\t\t%d\t\t%d\t\t%d\t\t%d\t\t%d\n",tr_inser,tr_find,tr_iter,tr_ins,tr_delet)
  runtime.ReadMemStats(&ms)
  fmt.Printf("mallocs %v , frees %v\n%v\n",ms.Mallocs,ms.Frees,ms.BySize)
}
