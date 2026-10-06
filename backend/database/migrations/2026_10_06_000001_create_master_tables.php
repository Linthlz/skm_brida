<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

/**
 * Data induk survei: jenis layanan, periode tahunan, pengaturan (satu baris),
 * dan bank pertanyaan.
 */
return new class extends Migration
{
    public function up(): void
    {
        Schema::create('layanan', function (Blueprint $table) {
            $table->id();
            $table->string('nama', 100)->unique();
            $table->boolean('aktif')->default(true);
            $table->unsignedSmallInteger('urutan')->default(0);
            $table->timestamps();
        });

        Schema::create('periode', function (Blueprint $table) {
            $table->id();
            $table->unsignedSmallInteger('tahun')->unique();
            // 4 = Permen PANRB 14/2017, 5 = skala yang dipakai survei ini.
            $table->unsignedTinyInteger('skala_maks')->default(5);
            $table->date('tanggal_tutup')->nullable();
            $table->timestamps();
        });

        Schema::create('pengaturan', function (Blueprint $table) {
            $table->id();
            $table->string('judul', 150);
            $table->foreignId('periode_aktif_id')->constrained('periode')->restrictOnDelete();
            $table->boolean('anonim')->default(true);
            $table->boolean('publik')->default(true);
            $table->timestamps();
        });

        Schema::create('pertanyaan', function (Blueprint $table) {
            $table->id();
            $table->string('kode', 10)->unique();
            $table->string('unsur', 100);
            $table->text('teks');
            $table->boolean('wajib')->default(true);
            $table->boolean('aktif')->default(true);
            $table->unsignedSmallInteger('urutan');
            $table->timestamps();
            // Soft delete: jawaban periode lama tetap merujuk ke pertanyaan ini.
            $table->softDeletes();

            $table->index(['aktif', 'urutan']);
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('pertanyaan');
        Schema::dropIfExists('pengaturan');
        Schema::dropIfExists('periode');
        Schema::dropIfExists('layanan');
    }
};
