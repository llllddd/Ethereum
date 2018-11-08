static int genkey(unsigned char *pk, unsigned char *sk){
    int ret_val;
    if ((ret_val = crypto_sign_keypair(pk,sk))!=0){
        return 1;
    }
    return 0;
}


static int getpk(unsigned char *sk, unsigned char *pk){
  int ret_val;
  if ((ret_val = crypto_get_pk(sk,pk))!=0){
    return 1;
  }
  return 0;
}

static int chain_sign(unsigned char *sm,
               unsigned long long smlen,
               const unsigned char *m,
               unsigned long long mlen,
               const unsigned char *sk)
{
  int val;
  if((val = crypto_sign(sm,&smlen,m,mlen,sk))!=0){
      return 1;
  }
    return 0;
}

static int chain_verify_sign(unsigned char *m,
                     unsigned long long *mlen,
                     const unsigned char *sm,
                     unsigned long long smlen,
                     const unsigned char *pk)
{
  int val;
  if ((val = crypto_sign_open(m,mlen,sm,smlen,pk))!=0){
     return 1;
  }
   return 0;
}
