<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Attributes\Fillable;
use Illuminate\Database\Eloquent\Attributes\Table;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

/** Saran dari responden beserta status tindak lanjutnya. */
#[Table('feedback')]
#[Fillable(['isi', 'kategori_id', 'status', 'catatan'])]
class Feedback extends Model
{
    /** Status hanya boleh maju satu langkah. */
    public const BERIKUT = [
        'Baru' => 'Ditinjau',
        'Ditinjau' => 'Ditindaklanjuti',
        'Ditindaklanjuti' => 'Selesai',
    ];

    public function responden(): BelongsTo
    {
        return $this->belongsTo(Responden::class);
    }

    public function kategori(): BelongsTo
    {
        return $this->belongsTo(Pertanyaan::class, 'kategori_id')->withTrashed();
    }

    public function penangan(): BelongsTo
    {
        return $this->belongsTo(User::class, 'ditangani_oleh');
    }

    /** Nomor tampilan, contoh FB-0241. Diturunkan dari id agar tidak perlu kolom terpisah. */
    public function nomor(): string
    {
        return 'FB-'.str_pad((string) $this->id, 4, '0', STR_PAD_LEFT);
    }
}
