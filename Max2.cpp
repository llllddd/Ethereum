#include<iostream>
#include<algorithm>


void Max2(int* a, int lo, int hi,int &x1,int &x2){
  x1=lo;
  x2 = lo;
  for(int i=lo+1;i<hi;i++){
    if(a[x1]<a[i]) x1=i;
  }
  for(int j=lo+1;j<x1;j++){
    if(a[x2]<a[j])x2=j;
  }
  for(int k= x1+1 ;k<hi;k++){
    if(a[x2]<a[k])x2=k;
  }

}

void Max2_a(int* a, int lo, int hi,int &x1,int &x2){
  x1=lo;
  x2=lo+1;
  if(a[x1]<a[x2]) std::swap(x1,x2);
  for (int i=lo+2;i<hi;i++){
    if (a[x2]<a[i]){
         x2=i;
      if(a[x1]<a[x2])
	std::swap(x1,x2);
    }
  }
  
}

void Max2_b(int* a, int lo, int hi,int &x1,int &x2){
  if(lo+2==hi){if(a[lo]<a[lo+1]){x1=lo+1;x2=lo;}x1=lo;x2=lo+1;}
  if(lo+3==hi){x1=lo;x2=lo+1;
    if(a[x1]<a[x2])
      std::swap(x1,x2);
    if(a[x1]<a[lo+3])
      x1=lo+3;
    if(a[x2]<a[lo+3])
      x2=lo+3;}
  int mid = (lo+hi)>>1;
  int x1_L;int x2_L; Max2_b(a,lo,mid,x1_L,x2_L);
  int x1_R;int x2_R; Max2_b(a,mid,hi,x1_R,x2_R);
  if(a[x1_L] < a[x1_R]){
    x1=x1_R;
    (a[x2_R]<a[x1_L])?x2=x1_L:x2=x2_R;

  }
  else{
     x1=x1_L;
(a[x2_L]<a[x1_R])?x2=x1_R:x2=x2_L;
    }

}


int main(){
  int a[]={1,2,3,4,5,6,7,8,9};
  int x1,x2;
  Max2_b(a,0,8,x1,x2);
  std::cout<<a[x1]<<a[x2]<<std::endl;
}

