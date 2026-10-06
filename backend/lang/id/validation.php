<?php

/*
 * Pesan validasi berbahasa Indonesia untuk aturan yang dipakai API ini.
 * Aturan lain jatuh ke bahasa cadangan (en).
 */
return [
    'accepted' => 'Centang persetujuan untuk mengirim jawaban Anda.',
    'array' => ':Attribute harus berupa daftar.',
    'between' => [
        'numeric' => ':Attribute harus antara :min dan :max.',
        'string' => ':Attribute harus antara :min dan :max karakter.',
        'array' => ':Attribute harus berisi :min sampai :max item.',
    ],
    'boolean' => ':Attribute harus bernilai benar atau salah.',
    'date' => ':Attribute bukan tanggal yang valid.',
    'date_format' => ':Attribute harus berformat :format.',
    'distinct' => ':Attribute memiliki nilai duplikat.',
    'email' => 'Masukkan alamat email yang valid, contoh: nama@email.com',
    'exists' => ':Attribute yang dipilih tidak valid.',
    'in' => ':Attribute yang dipilih tidak valid.',
    'integer' => ':Attribute harus berupa bilangan bulat.',
    'max' => [
        'numeric' => ':Attribute tidak boleh lebih dari :max.',
        'string' => ':Attribute tidak boleh lebih dari :max karakter.',
        'array' => ':Attribute tidak boleh lebih dari :max item.',
    ],
    'min' => [
        'numeric' => ':Attribute minimal :min.',
        'string' => ':Attribute minimal :min karakter.',
        'array' => ':Attribute minimal berisi :min item.',
    ],
    'prohibited' => 'Kolom :attribute tidak boleh diisi.',
    'required' => ':Attribute wajib diisi.',
    'string' => ':Attribute harus berupa teks.',

    'custom' => [
        'jk' => ['required' => 'Silakan pilih salah satu jawaban sebelum melanjutkan.'],
        'frekuensi' => ['required' => 'Silakan pilih salah satu jawaban sebelum melanjutkan.'],
        'usia' => ['between' => 'Masukkan usia antara 10 dan 99 tahun.'],
        'pendidikan' => ['required' => 'Silakan pilih pendidikan terakhir Anda.'],
        'pekerjaan' => ['required' => 'Silakan pilih pekerjaan Anda.'],
        'layanan_id' => ['required' => 'Silakan pilih layanan yang Anda gunakan.'],
        'topik' => ['required' => 'Silakan pilih salah satu topik.'],
        'jawaban.*.nilai' => ['between' => 'Nilai harus antara :min dan :max.'],
    ],

    'attributes' => [
        'nama' => 'nama',
        'email' => 'email',
        'password' => 'kata sandi',
        'jk' => 'jenis kelamin',
        'usia' => 'usia',
        'pendidikan' => 'pendidikan',
        'pekerjaan' => 'pekerjaan',
        'frekuensi' => 'frekuensi',
        'kecamatan' => 'kecamatan',
        'layanan_id' => 'layanan',
        'jawaban' => 'jawaban',
        'jawaban.*.pertanyaan_id' => 'pertanyaan',
        'jawaban.*.nilai' => 'nilai',
        'saran' => 'saran',
        'apresiasi' => 'apresiasi',
        'setuju' => 'persetujuan',
        'topik' => 'topik',
        'pesan' => 'pesan',
        'unsur' => 'unsur pelayanan',
        'teks' => 'teks pertanyaan',
        'judul' => 'judul survei',
        'periode' => 'periode',
        'skala_maks' => 'skala penilaian',
        'tanggal_tutup' => 'tanggal penutupan',
        'catatan' => 'catatan tindak lanjut',
        'kategori_id' => 'kategori',
        'status' => 'status',
        'per_page' => 'jumlah per halaman',
        'website' => 'website',
    ],
];
