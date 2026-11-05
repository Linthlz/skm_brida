package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/itsXanadium/singa-riset-backend/config"
	"github.com/itsXanadium/singa-riset-backend/controllers"
	"github.com/itsXanadium/singa-riset-backend/middleware"
	"github.com/itsXanadium/singa-riset-backend/services"
	"github.com/itsXanadium/singa-riset-backend/utils"
)

func skmLimiter(max int, expiration time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: expiration,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false, "status": "Too Many Requests", "ResponseCode": fiber.StatusTooManyRequests,
				"message": "Terlalu banyak permintaan. Coba lagi beberapa saat lagi.", "error": "rate limit exceeded",
			})
		},
	})
}

// SetupSkm mendaftarkan endpoint Survei Kepuasan Masyarakat (SKM).
// Dua peran: responden (publik, tanpa akun) dan admin (users.role = "admin").
// Login admin memakai /api/v1/auth/login milik Singa Riset.
func SetupSkm(app *fiber.App, skm *controllers.SkmController, userService services.UserService) {
	//============ Public Endpoints (responden) ============
	publik := app.Group("/api/v2/skm")
	baca := skmLimiter(60, time.Minute)
	publik.Get("/referensi", baca, skm.Referensi)
	publik.Get("/survei/form", baca, skm.Form)
	publik.Get("/periode", baca, skm.Periode)
	publik.Get("/hasil", baca, skm.HasilPublik)
	publik.Get("/ringkasan", baca, skm.Ringkasan)
	// Survei bisa diisi bergantian dari satu perangkat/IP kantor (loket layanan),
	// jadi batasnya cukup longgar (SKM_SURVEI_PER_JAM di .env).
	publik.Post("/survei", skmLimiter(utils.SkmBatasPerJam("SKM_SURVEI_PER_JAM", 20), time.Hour), skm.KirimSurvei)
	publik.Post("/kontak", skmLimiter(utils.SkmBatasPerJam("SKM_KONTAK_PER_JAM", 5), time.Hour), skm.KirimKontak)

	//============ Admin Endpoints ===============
	admin := app.Group("/api/v1/admin/skm",
		jwtware.New(jwtware.Config{
			SigningKey: []byte(config.AppConfig.JWTSecret),
			ContextKey: "user",
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				return utils.UnauthorizedAccess(c, "Unauthorized", "missing, invalid or expired token")
			},
		}),
		middleware.Verified(userService),
		middleware.AdminOnly(userService),
	)
	admin.Get("/dashboard", skm.Dashboard)
	admin.Get("/hasil", skm.Hasil)
	admin.Get("/statistik", skm.Statistik)

	admin.Get("/responden", skm.Responden)
	admin.Get("/responden/:id", skm.DetailResponden)

	admin.Get("/pertanyaan", skm.Pertanyaan)
	admin.Post("/pertanyaan", skm.TambahPertanyaan)
	admin.Put("/pertanyaan/urutan", skm.UrutanPertanyaan)
	admin.Patch("/pertanyaan/:id", skm.UbahPertanyaan)
	admin.Delete("/pertanyaan/:id", skm.HapusPertanyaan)

	admin.Get("/feedback", skm.Feedback)
	admin.Get("/feedback/hitung", skm.HitungFeedback)
	admin.Get("/feedback/:id", skm.DetailFeedback)
	admin.Patch("/feedback/:id", skm.UbahFeedback)

	admin.Get("/pesan-kontak", skm.PesanKontak)
	admin.Get("/pesan-kontak/hitung", skm.HitungPesan)
	admin.Patch("/pesan-kontak/:id", skm.UbahPesan)

	admin.Get("/laporan/export", skm.Laporan)

	admin.Get("/pengaturan", skm.Pengaturan)
	admin.Put("/pengaturan", skm.SimpanPengaturan)
}
