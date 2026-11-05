package controllers

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/itsXanadium/singa-riset-backend/models"
	"github.com/itsXanadium/singa-riset-backend/repositories"
	"github.com/itsXanadium/singa-riset-backend/services"
	"github.com/itsXanadium/singa-riset-backend/utils"
)

// SkmController: Survei Kepuasan Masyarakat. Endpoint publik (/api/v2/skm)
// dipakai responden tanpa akun; endpoint admin (/api/v1/admin/skm) hanya
// untuk users dengan role admin.
type SkmController struct {
	service services.SkmService
}

func NewSkmController(service services.SkmService) *SkmController {
	return &SkmController{service: service}
}

const skmPesanValidasi = "Validasi gagal"

func skmID(ctx *fiber.Ctx) (int64, bool) {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	return id, err == nil && id > 0
}

func skmQueryInt(ctx *fiber.Ctx, key string) (int64, bool) {
	v := strings.TrimSpace(ctx.Query(key))
	if v == "" {
		return 0, true
	}
	n, err := strconv.ParseInt(v, 10, 64)
	return n, err == nil && n > 0
}

// skmPeriodeFilter membaca filter bersama: periode (tahun), layanan_id, pendidikan.
func (c *SkmController) skmPeriodeFilter(ctx *fiber.Ctx) (*models.SkmPeriode, repositories.SkmFilter, error) {
	tahun, ok := skmQueryInt(ctx, "periode")
	if !ok {
		return nil, repositories.SkmFilter{}, utils.ErrValidation("periode: Periode tidak valid.")
	}
	layananID, ok := skmQueryInt(ctx, "layanan_id")
	if !ok {
		return nil, repositories.SkmFilter{}, utils.ErrValidation("layanan_id: Layanan tidak valid.")
	}
	periode, err := c.service.Periode(int(tahun))
	if err != nil {
		return nil, repositories.SkmFilter{}, err
	}
	f, err := c.service.Filter(layananID, strings.TrimSpace(ctx.Query("pendidikan")))
	return periode, f, err
}

// ============ Publik ============

// Referensi: GET /api/v2/skm/referensi — pilihan isian form survei, kontak, dan filter.
func (c *SkmController) Referensi(ctx *fiber.Ctx) error {
	data, err := c.service.Referensi()
	if err != nil {
		return utils.InternalServerError(ctx, "Gagal mengambil data referensi", err.Error())
	}
	return utils.Success(ctx, "Data berhasil diambil", data)
}

// Form: GET /api/v2/skm/survei/form
func (c *SkmController) Form(ctx *fiber.Ctx) error {
	data, err := c.service.Form()
	if err != nil {
		return utils.InternalServerError(ctx, "Gagal mengambil form survei", err.Error())
	}
	return utils.Success(ctx, "Data berhasil diambil", data)
}

// KirimSurvei: POST /api/v2/skm/survei
func (c *SkmController) KirimSurvei(ctx *fiber.Ctx) error {
	var body services.SkmKirimSurvei
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Invalid request body", "failed to parse request body")
	}
	data, err := c.service.Kirim(body, ctx.IP())
	if err != nil {
		return utils.HandleError(ctx, skmPesanValidasi, err)
	}
	return utils.CreatedWithData(ctx, "Jawaban survei berhasil dikirim", data)
}

// KirimKontak: POST /api/v2/skm/kontak
func (c *SkmController) KirimKontak(ctx *fiber.Ctx) error {
	var body services.SkmKirimKontak
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Invalid request body", "failed to parse request body")
	}
	if err := c.service.KirimKontak(body, ctx.IP()); err != nil {
		return utils.HandleError(ctx, skmPesanValidasi, err)
	}
	return utils.Created(ctx, "Pesan terkirim ke tim pelayanan BRIDA")
}

// Periode: GET /api/v2/skm/periode
func (c *SkmController) Periode(ctx *fiber.Ctx) error {
	data, err := c.service.DaftarPeriodePublik()
	if err != nil {
		return utils.InternalServerError(ctx, "Gagal mengambil periode", err.Error())
	}
	return utils.Success(ctx, "Data berhasil diambil", data)
}

// HasilPublik: GET /api/v2/skm/hasil — ditolak bila publikasi dimatikan.
func (c *SkmController) HasilPublik(ctx *fiber.Ctx) error {
	if err := c.service.PastikanPublik(); err != nil {
		return utils.HandleError(ctx, "Hasil survei belum dipublikasikan.", err)
	}
	return c.Hasil(ctx)
}

// Ringkasan: GET /api/v2/skm/ringkasan — angka ringkas untuk halaman depan.
func (c *SkmController) Ringkasan(ctx *fiber.Ctx) error {
	if err := c.service.PastikanPublik(); err != nil {
		return utils.HandleError(ctx, "Hasil survei belum dipublikasikan.", err)
	}
	data, err := c.service.Ringkasan()
	if err != nil {
		return utils.InternalServerError(ctx, "Gagal mengambil ringkasan", err.Error())
	}
	return utils.Success(ctx, "Data berhasil diambil", data)
}

// ============ Admin: hasil ============

// Hasil: GET /api/v1/admin/skm/hasil?periode=&layanan_id=&pendidikan=
func (c *SkmController) Hasil(ctx *fiber.Ctx) error {
	periode, f, err := c.skmPeriodeFilter(ctx)
	if err != nil {
		return utils.HandleError(ctx, skmPesanValidasi, err)
	}
	data, err := c.service.Rekap(periode, f)
	if err != nil {
		return utils.InternalServerError(ctx, "Gagal menghitung hasil survei", err.Error())
	}
	return utils.Success(ctx, "Data berhasil diambil", data)
}

// Dashboard: GET /api/v1/admin/skm/dashboard?periode=
func (c *SkmController) Dashboard(ctx *fiber.Ctx) error {
	periode, _, err := c.skmPeriodeFilter(ctx)
	if err != nil {
		return utils.HandleError(ctx, skmPesanValidasi, err)
	}
	data, err := c.service.Dashboard(periode)
	if err != nil {
		return utils.InternalServerError(ctx, "Gagal mengambil dashboard", err.Error())
	}
	return utils.Success(ctx, "Data berhasil diambil", data)
}

// Statistik: GET /api/v1/admin/skm/statistik?periode=
func (c *SkmController) Statistik(ctx *fiber.Ctx) error {
	periode, _, err := c.skmPeriodeFilter(ctx)
	if err != nil {
		return utils.HandleError(ctx, skmPesanValidasi, err)
	}
	data, err := c.service.Statistik(periode)
	if err != nil {
		return utils.InternalServerError(ctx, "Gagal mengambil statistik", err.Error())
	}
	return utils.Success(ctx, "Data berhasil diambil", data)
}

// Laporan: GET /api/v1/admin/skm/laporan/export?format=xlsx|pdf&periode=&layanan_id=&pendidikan=
func (c *SkmController) Laporan(ctx *fiber.Ctx) error {
	periode, f, err := c.skmPeriodeFilter(ctx)
	if err != nil {
		return utils.HandleError(ctx, skmPesanValidasi, err)
	}
	berkas, err := c.service.Laporan(periode, f, ctx.Query("format"))
	if err != nil {
		return utils.HandleError(ctx, "Gagal membuat laporan", err)
	}
	ctx.Set(fiber.HeaderContentType, berkas.ContentType)
	ctx.Set(fiber.HeaderContentDisposition, `attachment; filename="`+berkas.Nama+`"`)
	return ctx.Send(berkas.Isi)
}

// ============ Admin: responden ============

// Responden: GET /api/v1/admin/skm/responden?periode=&layanan_id=&search=&sort=&dir=&page=&limit=
func (c *SkmController) Responden(ctx *fiber.Ctx) error {
	page, limit, offset := paging(ctx, 10)
	tahun, ok := skmQueryInt(ctx, "periode")
	if !ok {
		return utils.BadReq(ctx, skmPesanValidasi, "periode: Periode tidak valid.")
	}
	layananID, ok := skmQueryInt(ctx, "layanan_id")
	if !ok {
		return utils.BadReq(ctx, skmPesanValidasi, "layanan_id: Layanan tidak valid.")
	}
	search := strings.TrimSpace(ctx.Query("search"))
	sort := ctx.Query("sort")
	list, total, err := c.service.DaftarResponden(services.SkmRespondenParam{
		Tahun: int(tahun), LayananID: layananID, Search: search, Sort: sort, Dir: ctx.Query("dir"),
		Limit: limit, Offset: offset,
	})
	if err != nil {
		return utils.HandleError(ctx, skmPesanValidasi, err)
	}
	return utils.SuccessPaginated(ctx, "Data berhasil diambil", list, metaFor(page, limit, total, search, sort))
}

// DetailResponden: GET /api/v1/admin/skm/responden/:id
func (c *SkmController) DetailResponden(ctx *fiber.Ctx) error {
	id, ok := skmID(ctx)
	if !ok {
		return utils.BadReq(ctx, "Invalid ID", "id must be a positive integer")
	}
	data, err := c.service.DetailResponden(id)
	if err != nil {
		return utils.HandleError(ctx, "Data tidak ditemukan", err)
	}
	return utils.Success(ctx, "Data berhasil ditemukan", data)
}

// ============ Admin: bank pertanyaan ============

// Pertanyaan: GET /api/v1/admin/skm/pertanyaan
func (c *SkmController) Pertanyaan(ctx *fiber.Ctx) error {
	data, err := c.service.DaftarPertanyaan()
	if err != nil {
		return utils.InternalServerError(ctx, "Gagal mengambil pertanyaan", err.Error())
	}
	return utils.Success(ctx, "Data berhasil diambil", data)
}

// TambahPertanyaan: POST /api/v1/admin/skm/pertanyaan
func (c *SkmController) TambahPertanyaan(ctx *fiber.Ctx) error {
	var body services.SkmPertanyaanInput
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Invalid request body", "failed to parse request body")
	}
	data, err := c.service.TambahPertanyaan(body)
	if err != nil {
		return utils.HandleError(ctx, skmPesanValidasi, err)
	}
	return utils.CreatedWithData(ctx, "Pertanyaan ditambahkan", data)
}

// UbahPertanyaan: PATCH /api/v1/admin/skm/pertanyaan/:id
func (c *SkmController) UbahPertanyaan(ctx *fiber.Ctx) error {
	id, ok := skmID(ctx)
	if !ok {
		return utils.BadReq(ctx, "Invalid ID", "id must be a positive integer")
	}
	var body services.SkmPertanyaanInput
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Invalid request body", "failed to parse request body")
	}
	data, err := c.service.UbahPertanyaan(id, body)
	if err != nil {
		return utils.HandleError(ctx, skmPesanValidasi, err)
	}
	return utils.Success(ctx, "Pertanyaan diperbarui", data)
}

// HapusPertanyaan: DELETE /api/v1/admin/skm/pertanyaan/:id
func (c *SkmController) HapusPertanyaan(ctx *fiber.Ctx) error {
	id, ok := skmID(ctx)
	if !ok {
		return utils.BadReq(ctx, "Invalid ID", "id must be a positive integer")
	}
	if err := c.service.HapusPertanyaan(id); err != nil {
		return utils.HandleError(ctx, "Gagal menghapus pertanyaan", err)
	}
	return utils.Success(ctx, "Pertanyaan dihapus", nil)
}

// UrutanPertanyaan: PUT /api/v1/admin/skm/pertanyaan/urutan
func (c *SkmController) UrutanPertanyaan(ctx *fiber.Ctx) error {
	var body struct {
		IDs []int64 `json:"ids"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Invalid request body", "failed to parse request body")
	}
	data, err := c.service.UrutkanPertanyaan(body.IDs)
	if err != nil {
		return utils.HandleError(ctx, skmPesanValidasi, err)
	}
	return utils.Success(ctx, "Urutan pertanyaan disimpan", data)
}

// ============ Admin: feedback & pesan kontak ============

// Feedback: GET /api/v1/admin/skm/feedback?status=&search=&page=&limit=
func (c *SkmController) Feedback(ctx *fiber.Ctx) error {
	page, limit, offset := paging(ctx, 10)
	search := strings.TrimSpace(ctx.Query("search"))
	list, total, err := c.service.DaftarFeedback(ctx.Query("status"), search, limit, offset)
	if err != nil {
		return utils.HandleError(ctx, skmPesanValidasi, err)
	}
	return utils.SuccessPaginated(ctx, "Data berhasil diambil", list, metaFor(page, limit, total, search, ""))
}

// HitungFeedback: GET /api/v1/admin/skm/feedback/hitung — jumlah per status.
func (c *SkmController) HitungFeedback(ctx *fiber.Ctx) error {
	data, err := c.service.HitungStatusFeedback()
	if err != nil {
		return utils.InternalServerError(ctx, "Gagal menghitung feedback", err.Error())
	}
	return utils.Success(ctx, "Data berhasil diambil", data)
}

// DetailFeedback: GET /api/v1/admin/skm/feedback/:id
func (c *SkmController) DetailFeedback(ctx *fiber.Ctx) error {
	id, ok := skmID(ctx)
	if !ok {
		return utils.BadReq(ctx, "Invalid ID", "id must be a positive integer")
	}
	data, err := c.service.DetailFeedback(id)
	if err != nil {
		return utils.HandleError(ctx, "Data tidak ditemukan", err)
	}
	return utils.Success(ctx, "Data berhasil ditemukan", data)
}

// UbahFeedback: PATCH /api/v1/admin/skm/feedback/:id — body {status?, catatan?, kategori_id?}
func (c *SkmController) UbahFeedback(ctx *fiber.Ctx) error {
	actorID, ok := currentUserID(ctx)
	if !ok {
		return utils.UnauthorizedAccess(ctx, "Unauthorized", "invalid token")
	}
	id, ok := skmID(ctx)
	if !ok {
		return utils.BadReq(ctx, "Invalid ID", "id must be a positive integer")
	}
	var body services.SkmFeedbackInput
	var ada map[string]json.RawMessage
	if json.Unmarshal(ctx.Body(), &body) != nil || json.Unmarshal(ctx.Body(), &ada) != nil {
		return utils.BadReq(ctx, "Invalid request body", "failed to parse request body")
	}
	_, body.CatatanAda = ada["catatan"]
	_, body.KategoriIDAda = ada["kategori_id"]
	data, err := c.service.UbahFeedback(id, body, actorID)
	if err != nil {
		return utils.HandleError(ctx, skmPesanValidasi, err)
	}
	return utils.Success(ctx, "Feedback diperbarui", data)
}

// PesanKontak: GET /api/v1/admin/skm/pesan-kontak?dibaca=&page=&limit=
func (c *SkmController) PesanKontak(ctx *fiber.Ctx) error {
	page, limit, offset := paging(ctx, 10)
	var dibaca *bool
	if v := ctx.Query("dibaca"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return utils.BadReq(ctx, skmPesanValidasi, "dibaca: Nilai harus true atau false.")
		}
		dibaca = &b
	}
	list, total, err := c.service.DaftarPesan(dibaca, limit, offset)
	if err != nil {
		return utils.InternalServerError(ctx, "Gagal mengambil pesan", err.Error())
	}
	return utils.SuccessPaginated(ctx, "Data berhasil diambil", list, metaFor(page, limit, total, "", ""))
}

// HitungPesan: GET /api/v1/admin/skm/pesan-kontak/hitung
func (c *SkmController) HitungPesan(ctx *fiber.Ctx) error {
	n, err := c.service.JumlahPesanBelumDibaca()
	if err != nil {
		return utils.InternalServerError(ctx, "Gagal menghitung pesan", err.Error())
	}
	return utils.Success(ctx, "Data berhasil diambil", fiber.Map{"belum_dibaca": n})
}

// UbahPesan: PATCH /api/v1/admin/skm/pesan-kontak/:id — body {dibaca}
func (c *SkmController) UbahPesan(ctx *fiber.Ctx) error {
	id, ok := skmID(ctx)
	if !ok {
		return utils.BadReq(ctx, "Invalid ID", "id must be a positive integer")
	}
	var body struct {
		Dibaca *bool `json:"dibaca"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Invalid request body", "failed to parse request body")
	}
	if body.Dibaca == nil {
		return utils.BadReq(ctx, skmPesanValidasi, "dibaca: Wajib diisi.")
	}
	data, err := c.service.UbahPesan(id, *body.Dibaca)
	if err != nil {
		return utils.HandleError(ctx, "Gagal memperbarui pesan", err)
	}
	return utils.Success(ctx, "Pesan diperbarui", data)
}

// ============ Admin: pengaturan ============

// Pengaturan: GET /api/v1/admin/skm/pengaturan
func (c *SkmController) Pengaturan(ctx *fiber.Ctx) error {
	data, err := c.service.Pengaturan()
	if err != nil {
		return utils.InternalServerError(ctx, "Gagal mengambil pengaturan", err.Error())
	}
	return utils.Success(ctx, "Data berhasil diambil", data)
}

// SimpanPengaturan: PUT /api/v1/admin/skm/pengaturan
func (c *SkmController) SimpanPengaturan(ctx *fiber.Ctx) error {
	var body services.SkmPengaturanInput
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Invalid request body", "failed to parse request body")
	}
	data, err := c.service.SimpanPengaturan(body)
	if err != nil {
		return utils.HandleError(ctx, skmPesanValidasi, err)
	}
	return utils.Success(ctx, "Pengaturan disimpan", data)
}
