Edit `account` dan `key` menggunakan alamat email cloudflare dan Global API Keys.
```plaintext
    account := "test@domain.com"
    key := "123"
```

List IP Access Rules.
```sh
go run cloudflare.go
```

Menambahkan IP ke daftar `block`.
```sh
go run cloudflare.go --mode block --ip 198.51.100.4 --notes "bruteforce"
```

Selain `block`, Anda juga dapat menggunakan mode `challenge`, `whitelist`, `js_challenge`, `managed_challenge`.
Cek dokumentasi [Create an IP Access rule](https://developers.cloudflare.com/api/operations/ip-access-rules-for-a-user-create-an-ip-access-rule).

Menghapus IP dari daftar `block`.
```sh
go run cloudflare.go --del --mode block --ip 198.51.100.4
```
