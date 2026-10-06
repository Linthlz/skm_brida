<?php

namespace App\Services;

use App\Models\Feedback;
use App\Models\User;
use Illuminate\Validation\ValidationException;

class FeedbackService
{
    /** Hitungan per status untuk kartu filter di halaman Feedback. */
    public function hitungStatus(): array
    {
        $hitung = Feedback::selectRaw('status, COUNT(*) as n')->groupBy('status')->pluck('n', 'status');

        return collect(config('skm.status_feedback'))->mapWithKeys(fn ($s) => [$s => (int) ($hitung[$s] ?? 0)])->all();
    }

    /** Status hanya boleh maju satu langkah (Baru → Ditinjau → Ditindaklanjuti → Selesai). */
    public function perbarui(Feedback $feedback, array $data, User $oleh): Feedback
    {
        if (isset($data['status']) && $data['status'] !== $feedback->status) {
            $berikut = Feedback::BERIKUT[$feedback->status] ?? null;
            if ($data['status'] !== $berikut) {
                throw ValidationException::withMessages([
                    'status' => $berikut
                        ? "Status {$feedback->status} hanya dapat diubah menjadi {$berikut}."
                        : 'Feedback yang sudah selesai tidak dapat diubah statusnya.',
                ]);
            }
        }

        $feedback->fill($data);
        $feedback->ditangani_oleh = $oleh->id;
        $feedback->save();

        return $feedback->load(['responden.layanan', 'kategori', 'penangan']);
    }
}
