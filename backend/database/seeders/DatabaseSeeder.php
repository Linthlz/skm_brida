<?php

namespace Database\Seeders;

use Illuminate\Database\Seeder;

class DatabaseSeeder extends Seeder
{
    public function run(): void
    {
        $this->call([MasterSeeder::class, AdminSeeder::class]);

        // Data contoh hanya untuk pengembangan, tidak pernah di production.
        if (app()->environment('local')) {
            $this->call(DataContohSeeder::class);
        }
    }
}
