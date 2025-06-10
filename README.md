DEBIT    = Mengurangi wallet
CREDIT   = Menambag Wallet


#_FLOW :
user melakukan transaksi dengan tipe
-TOPUP      : menambah wallet
-PURCHASE   : mengurangi wallet
-REFUND     : mengembalikan nilai wallet


#_STATUS Transaksi:
-PENDING    : Menjadi nilai default
-SUCCESS    : Akan melakukan balance (DEBIT atau CREDIT)
-FAILED     : Transaksi di Gagalkan
-REVERSED   : Pembatalan transaksi (status success) dengan durasi bisa di rubah darai transaski dibuat