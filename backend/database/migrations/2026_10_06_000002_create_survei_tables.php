<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

/**
 * Hasil pengisian survei: identitas responden, jawaban per pertanyaan,
 * saran (feedback), dan pesan dari form kontak.
 */
return new class extends Migration
{
    public function up(): void
    {
        Schema::create('responden', function (Blueprint $table) {
            $table->id();
            $table->foreignId('periode_id')->constrained('periode')->restrictOnDelete();
            // Nomor urut per periode; nomor tampilan SKM-<tahun>-<urut>.
            $table->unsignedInteger('urut');
            $table->string('nomor', 20)->unique();
            $table->foreignId('layanan_id')->constrained('layanan')->restrictOnDelete();
            $table->string('nama', 100)->nullable();
            $table->enum('jk', ['Laki-laki', 'Perempuan']);
            $table->unsignedTinyInteger('usia');
            $table->string('pendidikan', 40);
            $table->string('pekerjaan', 40);
            $table->string('frekuensi', 40);
            $table->string('kecamatan', 40)->nullable();
            $table->string('apresiasi', 300)->nullable();
            $table->decimal('rata_rata', 4, 2);
            $table->char('ip_hash', 64)->nullable();
            $table->timestamps();

            $table->unique(['periode_id', 'urut']);
            $table->index(['periode_id', 'created_at']);
            $table->index(['periode_id', 'layanan_id']);
            $table->index(['periode_id', 'pendidikan']);
            $table->index(['periode_id', 'pekerjaan']);
            $table->index(['ip_hash', 'created_at']);
        });

        Schema::create('jawaban', function (Blueprint $table) {
            $table->id();
            $table->foreignId('responden_id')->constrained('responden')->cascadeOnDelete();
            $table->foreignId('pertanyaan_id')->constrained('pertanyaan')->restrictOnDelete();
            $table->unsignedTinyInteger('nilai');

            $table->unique(['responden_id', 'pertanyaan_id']);
            $table->index(['pertanyaan_id', 'nilai']);
        });

        Schema::create('feedback', function (Blueprint $table) {
            $table->id();
            $table->foreignId('responden_id')->unique()->constrained('responden')->cascadeOnDelete();
            $table->foreignId('kategori_id')->nullable()->constrained('pertanyaan')->nullOnDelete();
            $table->string('isi', 600);
            $table->enum('status', ['Baru', 'Ditinjau', 'Ditindaklanjuti', 'Selesai'])->default('Baru')->index();
            $table->text('catatan')->nullable();
            $table->foreignId('ditangani_oleh')->nullable()->constrained('users')->nullOnDelete();
            $table->timestamps();

            $table->index('created_at');
        });

        Schema::create('pesan_kontak', function (Blueprint $table) {
            $table->id();
            $table->string('nama', 100);
            $table->string('email', 150);
            $table->string('topik', 40);
            $table->text('pesan');
            $table->boolean('dibaca')->default(false);
            $table->char('ip_hash', 64)->nullable();
            $table->timestamps();

            $table->index(['dibaca', 'created_at']);
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('pesan_kontak');
        Schema::dropIfExists('feedback');
        Schema::dropIfExists('jawaban');
        Schema::dropIfExists('responden');
    }
};
