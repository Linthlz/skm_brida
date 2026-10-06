<?php

/*
 * Daftar pilihan dan aturan domain SKM. Disamakan dengan src/data/opsi.js,
 * src/data/skala.js, dan src/utils/mutu.js di frontend.
 */
return [

    'jk' => ['Laki-laki', 'Perempuan'],

    'pendidikan' => [
        'SD / Sederajat', 'SMP / Sederajat', 'SMA / Sederajat',
        'Diploma (D1–D3)', 'Sarjana (S1/D4)', 'Pascasarjana (S2/S3)',
    ],

    'pekerjaan' => [
        'Pelajar / Mahasiswa', 'ASN / TNI / Polri', 'Pegawai Swasta', 'Wiraswasta',
        'Peneliti / Akademisi', 'Petani / Nelayan', 'Lainnya',
    ],

    'frekuensi' => ['Baru pertama kali', '2–3 kali', '4–6 kali', 'Lebih dari 6 kali'],

    'kecamatan' => [
        'Buleleng', 'Sukasada', 'Sawan', 'Banjar', 'Seririt',
        'Kubutambahan', 'Tejakula', 'Busungbiu', 'Gerokgak',
    ],

    'topik_kontak' => [
        'Pertanyaan tentang survei', 'Permintaan data hasil SKM', 'Pengaduan layanan',
        'Kerja sama riset & inovasi', 'Lainnya',
    ],

    'status_feedback' => ['Baru', 'Ditinjau', 'Ditindaklanjuti', 'Selesai'],

    // Label skala kepuasan per nilai, dikelompokkan menurut nilai maksimum skala.
    'label_skala' => [
        5 => [1 => 'Sangat Tidak Puas', 2 => 'Tidak Puas', 3 => 'Kurang Puas', 4 => 'Puas', 5 => 'Sangat Puas'],
        4 => [1 => 'Tidak Puas', 2 => 'Kurang Puas', 3 => 'Puas', 4 => 'Sangat Puas'],
    ],

    // Ambang mutu pelayanan (Permen PANRB 14/2017), batas bawah inklusif.
    'mutu' => [
        ['huruf' => 'A', 'min' => 88.31, 'label' => 'Sangat Baik'],
        ['huruf' => 'B', 'min' => 76.61, 'label' => 'Baik'],
        ['huruf' => 'C', 'min' => 65.00, 'label' => 'Kurang Baik'],
        ['huruf' => 'D', 'min' => 0, 'label' => 'Tidak Baik'],
    ],

    // Batas pengiriman per IP per jam untuk endpoint publik yang bisa disalahgunakan.
    'batas' => [
        'survei_per_jam' => (int) env('SKM_SURVEI_PER_JAM', 20),
        'kontak_per_jam' => (int) env('SKM_KONTAK_PER_JAM', 5),
    ],

    // Salt untuk hash IP responden/pengirim pesan (IP mentah tidak disimpan).
    'ip_salt' => env('SKM_IP_SALT', env('APP_KEY')),

    // Akun admin awal yang dibuat seeder. Tidak ada nilai bawaan agar
    // kata sandi tidak pernah tertulis di kode sumber.
    'admin' => [
        'nama' => env('ADMIN_NAMA', 'Administrator SKM'),
        'email' => env('ADMIN_EMAIL'),
        'password' => env('ADMIN_PASSWORD'),
        'jabatan' => env('ADMIN_JABATAN', 'Admin Pelayanan'),
    ],
];
