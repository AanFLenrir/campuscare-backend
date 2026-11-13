"""
CAMPUSCARE Backend — PDF Report Generator
Menggunakan ReportLab untuk membuat laporan profesional berbahasa Indonesia.
"""

from reportlab.lib.pagesizes import A4
from reportlab.lib.styles import getSampleStyleSheet, ParagraphStyle
from reportlab.lib.units import cm
from reportlab.lib import colors
from reportlab.lib.enums import TA_CENTER, TA_LEFT, TA_JUSTIFY
from reportlab.platypus import (
    SimpleDocTemplate, Paragraph, Spacer, Table, TableStyle,
    HRFlowable, PageBreak, KeepTogether
)
from reportlab.platypus.tableofcontents import TableOfContents
from reportlab.pdfbase import pdfmetrics
from reportlab.pdfbase.ttfonts import TTFont
import os

OUTPUT_PATH = os.path.join(os.path.dirname(__file__), "CAMPUSCARE-FINAL-REPORT.pdf")

# ── Color palette ─────────────────────────────────────────────────────────────
BLUE_DARK   = colors.HexColor("#1a237e")
BLUE_MID    = colors.HexColor("#283593")
BLUE_LIGHT  = colors.HexColor("#e8eaf6")
ACCENT      = colors.HexColor("#3949ab")
GRAY_DARK   = colors.HexColor("#37474f")
GRAY_MID    = colors.HexColor("#78909c")
GRAY_LIGHT  = colors.HexColor("#f5f5f5")
WHITE       = colors.white
BLACK       = colors.black
GREEN       = colors.HexColor("#2e7d32")
TABLE_HEAD  = colors.HexColor("#283593")
TABLE_ROW1  = colors.HexColor("#e8eaf6")
TABLE_ROW2  = colors.white

# ── Page template with header/footer ─────────────────────────────────────────
PAGE_WIDTH, PAGE_HEIGHT = A4

def make_doc(path):
    return SimpleDocTemplate(
        path,
        pagesize=A4,
        leftMargin=2.5*cm, rightMargin=2.5*cm,
        topMargin=2.8*cm, bottomMargin=2.5*cm,
        title="Laporan CAMPUSCARE Backend",
        author="MOHAMMAD ANHAZ ABDILAH",
        subject="Rancangan Sistem dan Boilerplate Backend CAMPUSCARE",
    )

def on_page(canvas, doc):
    """Draw header and footer on every page."""
    canvas.saveState()
    page_num = doc.page

    # Header bar
    canvas.setFillColor(BLUE_DARK)
    canvas.rect(doc.leftMargin, PAGE_HEIGHT - 1.8*cm,
                PAGE_WIDTH - doc.leftMargin - doc.rightMargin, 0.6*cm, fill=1, stroke=0)
    canvas.setFillColor(WHITE)
    canvas.setFont("Helvetica-Bold", 8)
    canvas.drawString(doc.leftMargin + 0.2*cm, PAGE_HEIGHT - 1.42*cm,
                      "CAMPUSCARE — Campus Facility Complaint and Management System")
    canvas.drawRightString(PAGE_WIDTH - doc.rightMargin - 0.2*cm,
                           PAGE_HEIGHT - 1.42*cm, "Universitas Airlangga")

    # Footer line
    canvas.setStrokeColor(BLUE_DARK)
    canvas.setLineWidth(0.5)
    canvas.line(doc.leftMargin, 1.8*cm,
                PAGE_WIDTH - doc.rightMargin, 1.8*cm)

    # Footer text
    canvas.setFillColor(GRAY_MID)
    canvas.setFont("Helvetica", 7.5)
    canvas.drawString(doc.leftMargin, 1.3*cm,
                      "D4 Teknik Informatika — Fakultas Vokasi — Universitas Airlangga")
    canvas.setFont("Helvetica-Bold", 8)
    canvas.drawRightString(PAGE_WIDTH - doc.rightMargin, 1.3*cm, f"Halaman {page_num}")

    canvas.restoreState()

# ── Styles ────────────────────────────────────────────────────────────────────
def build_styles():
    base = getSampleStyleSheet()

    styles = {
        "cover_title": ParagraphStyle(
            "cover_title", parent=base["Title"],
            fontSize=22, textColor=BLUE_DARK, spaceAfter=8,
            alignment=TA_CENTER, fontName="Helvetica-Bold", leading=28
        ),
        "cover_subtitle": ParagraphStyle(
            "cover_subtitle", parent=base["Normal"],
            fontSize=13, textColor=BLUE_MID, spaceAfter=6,
            alignment=TA_CENTER, fontName="Helvetica-Bold", leading=18
        ),
        "cover_body": ParagraphStyle(
            "cover_body", parent=base["Normal"],
            fontSize=11, textColor=GRAY_DARK, spaceAfter=4,
            alignment=TA_CENTER, fontName="Helvetica", leading=16
        ),
        "h1": ParagraphStyle(
            "h1", parent=base["Heading1"],
            fontSize=14, textColor=WHITE, spaceAfter=10, spaceBefore=18,
            fontName="Helvetica-Bold", leading=18, backColor=BLUE_DARK,
            leftIndent=-0.3*cm, rightIndent=-0.3*cm, borderPad=5
        ),
        "h2": ParagraphStyle(
            "h2", parent=base["Heading2"],
            fontSize=12, textColor=BLUE_DARK, spaceAfter=6, spaceBefore=12,
            fontName="Helvetica-Bold", leading=16
        ),
        "h3": ParagraphStyle(
            "h3", parent=base["Heading3"],
            fontSize=10.5, textColor=ACCENT, spaceAfter=4, spaceBefore=8,
            fontName="Helvetica-Bold", leading=14
        ),
        "body": ParagraphStyle(
            "body", parent=base["Normal"],
            fontSize=10, textColor=GRAY_DARK, spaceAfter=5,
            fontName="Helvetica", leading=15, alignment=TA_JUSTIFY
        ),
        "body_left": ParagraphStyle(
            "body_left", parent=base["Normal"],
            fontSize=10, textColor=GRAY_DARK, spaceAfter=4,
            fontName="Helvetica", leading=14, alignment=TA_LEFT
        ),
        "bullet": ParagraphStyle(
            "bullet", parent=base["Normal"],
            fontSize=10, textColor=GRAY_DARK, spaceAfter=3,
            fontName="Helvetica", leading=14, leftIndent=0.8*cm,
            bulletIndent=0.3*cm
        ),
        "code": ParagraphStyle(
            "code", parent=base["Code"],
            fontSize=8.5, textColor=GRAY_DARK, spaceAfter=3,
            fontName="Courier", leading=13, backColor=GRAY_LIGHT,
            leftIndent=0.5*cm, rightIndent=0.5*cm, borderPad=4
        ),
        "caption": ParagraphStyle(
            "caption", parent=base["Normal"],
            fontSize=9, textColor=GRAY_MID, spaceAfter=8,
            fontName="Helvetica-Oblique", leading=12, alignment=TA_CENTER
        ),
        "pass_label": ParagraphStyle(
            "pass_label", parent=base["Normal"],
            fontSize=10, textColor=GREEN, spaceAfter=2,
            fontName="Helvetica-Bold", leading=14, alignment=TA_CENTER
        ),
    }
    return styles

# ── Table helpers ─────────────────────────────────────────────────────────────
def make_table(headers, rows, col_widths=None, s=None):
    """Build a styled table."""
    header_cells = [Paragraph(f"<b>{h}</b>", ParagraphStyle(
        "th", fontSize=9, textColor=WHITE, fontName="Helvetica-Bold",
        leading=12, alignment=TA_CENTER
    )) for h in headers]

    data = [header_cells]
    for i, row in enumerate(rows):
        style = ParagraphStyle(
            "td", fontSize=9, textColor=GRAY_DARK, fontName="Helvetica",
            leading=12, alignment=TA_LEFT
        )
        data.append([Paragraph(str(cell), style) for cell in row])

    page_w = A4[0] - 5*cm
    if col_widths is None:
        n = len(headers)
        col_widths = [page_w / n] * n

    tbl = Table(data, colWidths=col_widths, repeatRows=1)

    ts = [
        ("BACKGROUND", (0, 0), (-1, 0), TABLE_HEAD),
        ("TEXTCOLOR",  (0, 0), (-1, 0), WHITE),
        ("FONTNAME",   (0, 0), (-1, 0), "Helvetica-Bold"),
        ("FONTSIZE",   (0, 0), (-1, 0), 9),
        ("ALIGN",      (0, 0), (-1, 0), "CENTER"),
        ("VALIGN",     (0, 0), (-1, -1), "MIDDLE"),
        ("GRID",       (0, 0), (-1, -1), 0.4, colors.HexColor("#b0bec5")),
        ("ROWBACKGROUNDS", (0, 1), (-1, -1), [TABLE_ROW1, TABLE_ROW2]),
        ("TOPPADDING",  (0, 0), (-1, -1), 5),
        ("BOTTOMPADDING", (0, 0), (-1, -1), 5),
        ("LEFTPADDING", (0, 0), (-1, -1), 6),
        ("RIGHTPADDING", (0, 0), (-1, -1), 6),
    ]
    tbl.setStyle(TableStyle(ts))
    return tbl

def info_box(label, value, s):
    """Small 2-col info box."""
    data = [[
        Paragraph(f"<b>{label}</b>", ParagraphStyle(
            "il", fontSize=9.5, textColor=WHITE, fontName="Helvetica-Bold", leading=13
        )),
        Paragraph(value, ParagraphStyle(
            "iv", fontSize=9.5, textColor=GRAY_DARK, fontName="Helvetica", leading=13
        ))
    ]]
    t = Table(data, colWidths=[4.5*cm, 11*cm])
    t.setStyle(TableStyle([
        ("BACKGROUND", (0, 0), (0, 0), BLUE_DARK),
        ("BACKGROUND", (1, 0), (1, 0), BLUE_LIGHT),
        ("VALIGN", (0, 0), (-1, -1), "MIDDLE"),
        ("TOPPADDING", (0, 0), (-1, -1), 5),
        ("BOTTOMPADDING", (0, 0), (-1, -1), 5),
        ("LEFTPADDING", (0, 0), (-1, -1), 8),
        ("GRID", (0, 0), (-1, -1), 0.3, colors.HexColor("#b0bec5")),
    ]))
    return t

def section_header(text, s):
    p = Paragraph(f"&nbsp;&nbsp;{text}", s["h1"])
    return p

# ── Cover Page ────────────────────────────────────────────────────────────────
def build_cover(s):
    story = []

    # Top decorative bar
    story.append(Spacer(1, 0.5*cm))
    bar_data = [[""]]
    bar = Table(bar_data, colWidths=[A4[0] - 5*cm], rowHeights=[0.8*cm])
    bar.setStyle(TableStyle([
        ("BACKGROUND", (0, 0), (-1, -1), BLUE_DARK),
        ("LINEBELOW", (0, 0), (-1, -1), 3, ACCENT),
    ]))
    story.append(bar)
    story.append(Spacer(1, 1.2*cm))

    # Title block
    story.append(Paragraph("LAPORAN RANCANGAN SISTEM", s["cover_title"]))
    story.append(Paragraph("DAN BOILERPLATE BACKEND", s["cover_title"]))
    story.append(Spacer(1, 0.4*cm))

    # Project name box
    proj_data = [["CAMPUSCARE"]]
    proj_tbl = Table(proj_data, colWidths=[A4[0] - 5*cm], rowHeights=[1.4*cm])
    proj_tbl.setStyle(TableStyle([
        ("BACKGROUND", (0, 0), (-1, -1), BLUE_MID),
        ("TEXTCOLOR",  (0, 0), (-1, -1), WHITE),
        ("FONTNAME",   (0, 0), (-1, -1), "Helvetica-Bold"),
        ("FONTSIZE",   (0, 0), (-1, -1), 20),
        ("ALIGN",      (0, 0), (-1, -1), "CENTER"),
        ("VALIGN",     (0, 0), (-1, -1), "MIDDLE"),
    ]))
    story.append(proj_tbl)
    story.append(Spacer(1, 0.3*cm))
    story.append(Paragraph(
        "Campus Facility Complaint and Management System",
        s["cover_subtitle"]
    ))
    story.append(Spacer(1, 1.5*cm))
    story.append(HRFlowable(width="100%", thickness=1, color=ACCENT))
    story.append(Spacer(1, 1.2*cm))

    # Identity section
    id_data = [
        ["Program Studi", "D4 Teknik Informatika"],
        ["Fakultas",      "Vokasi"],
        ["Universitas",   "Universitas Airlangga"],
        ["Tahun",         "2026"],
    ]
    id_tbl = Table(id_data, colWidths=[4.5*cm, 11*cm])
    id_tbl.setStyle(TableStyle([
        ("BACKGROUND", (0, 0), (0, -1), BLUE_DARK),
        ("BACKGROUND", (1, 0), (1, -1), BLUE_LIGHT),
        ("TEXTCOLOR",  (0, 0), (0, -1), WHITE),
        ("TEXTCOLOR",  (1, 0), (1, -1), GRAY_DARK),
        ("FONTNAME",   (0, 0), (-1, -1), "Helvetica-Bold"),
        ("FONTSIZE",   (0, 0), (-1, -1), 10.5),
        ("ALIGN",      (0, 0), (0, -1), "LEFT"),
        ("ALIGN",      (1, 0), (1, -1), "LEFT"),
        ("VALIGN",     (0, 0), (-1, -1), "MIDDLE"),
        ("TOPPADDING", (0, 0), (-1, -1), 7),
        ("BOTTOMPADDING", (0, 0), (-1, -1), 7),
        ("LEFTPADDING", (0, 0), (-1, -1), 10),
        ("GRID", (0, 0), (-1, -1), 0.3, colors.HexColor("#b0bec5")),
    ]))
    story.append(id_tbl)
    story.append(Spacer(1, 1.2*cm))
    story.append(HRFlowable(width="100%", thickness=1, color=ACCENT))
    story.append(Spacer(1, 0.8*cm))

    # Author
    story.append(Paragraph("Disusun Oleh:", s["cover_body"]))
    story.append(Spacer(1, 0.3*cm))
    author_data = [
        ["Nama",           "MOHAMMAD ANHAZ ABDILAH"],
        ["NIM",            "434241102"],
        ["Peran",          "Developer & Git Master"],
        ["GitHub",         "github.com/AanFLenrir"],
        ["Repository",     "github.com/AanFLenrir/campuscare-backend"],
    ]
    auth_tbl = Table(author_data, colWidths=[4.5*cm, 11*cm])
    auth_tbl.setStyle(TableStyle([
        ("BACKGROUND", (0, 0), (0, -1), ACCENT),
        ("BACKGROUND", (1, 0), (1, -1), colors.HexColor("#f3f4ff")),
        ("TEXTCOLOR",  (0, 0), (0, -1), WHITE),
        ("TEXTCOLOR",  (1, 0), (1, -1), GRAY_DARK),
        ("FONTNAME",   (0, 0), (-1, -1), "Helvetica-Bold"),
        ("FONTSIZE",   (0, 0), (-1, -1), 10),
        ("VALIGN",     (0, 0), (-1, -1), "MIDDLE"),
        ("TOPPADDING", (0, 0), (-1, -1), 6),
        ("BOTTOMPADDING", (0, 0), (-1, -1), 6),
        ("LEFTPADDING", (0, 0), (-1, -1), 10),
        ("GRID", (0, 0), (-1, -1), 0.3, colors.HexColor("#b0bec5")),
    ]))
    story.append(auth_tbl)
    story.append(Spacer(1, 1.5*cm))

    # Note for other members
    note_data = [[Paragraph(
        "<i>Catatan: Kolom anggota kelompok 2–4 harus dilengkapi dengan data asli "
        "sebelum laporan dikumpulkan. Jangan mengisi dengan identitas yang tidak sesuai.</i>",
        ParagraphStyle("note", fontSize=9, textColor=GRAY_MID, fontName="Helvetica-Oblique",
                       leading=13, alignment=TA_CENTER)
    )]]
    note_tbl = Table(note_data, colWidths=[A4[0] - 5*cm])
    note_tbl.setStyle(TableStyle([
        ("BACKGROUND", (0, 0), (-1, -1), colors.HexColor("#fff9c4")),
        ("LINEABOVE",  (0, 0), (-1, -1), 1, colors.HexColor("#f9a825")),
        ("LINEBELOW",  (0, 0), (-1, -1), 1, colors.HexColor("#f9a825")),
        ("TOPPADDING", (0, 0), (-1, -1), 8),
        ("BOTTOMPADDING", (0, 0), (-1, -1), 8),
        ("LEFTPADDING", (0, 0), (-1, -1), 10),
        ("RIGHTPADDING", (0, 0), (-1, -1), 10),
    ]))
    story.append(note_tbl)
    story.append(PageBreak())
    return story

# ── Build full document ───────────────────────────────────────────────────────
def build_story(s):
    story = build_cover(s)

    # ── 1. Latar Belakang ─────────────────────────────────────────────────────
    story.append(section_header("1. Latar Belakang", s))
    story.append(Spacer(1, 0.3*cm))
    story.append(Paragraph(
        "Fasilitas kampus seperti ruang kelas, laboratorium, perpustakaan, dan area umum "
        "sering mengalami kerusakan yang tidak segera dilaporkan maupun ditangani. "
        "Ketiadaan sistem pelaporan yang terstruktur menyebabkan kerusakan terakumulasi, "
        "mengganggu proses belajar mengajar, dan menurunkan kualitas layanan kampus.",
        s["body"]
    ))
    story.append(Paragraph(
        "Proses pelaporan manual melalui surat, telepon, atau datang langsung ke bagian "
        "pemeliharaan memiliki kelemahan seperti tidak tercatat dengan baik, sulit dipantau "
        "statusnya, dan tidak ada mekanisme akuntabilitas yang jelas antara pelapor, "
        "petugas, dan admin.",
        s["body"]
    ))
    story.append(Paragraph("<b>CAMPUSCARE</b> hadir sebagai solusi digital untuk:", s["body_left"]))
    for item in [
        "Memudahkan mahasiswa melaporkan kerusakan fasilitas secara sistematis dan terdigitalisasi.",
        "Membantu admin mengelola dan mendistribusikan pengaduan kepada petugas yang tepat.",
        "Memungkinkan petugas memperbarui status penanganan secara real-time.",
        "Menyediakan riwayat penanganan yang dapat diaudit dan dipertanggungjawabkan.",
        "Menghasilkan laporan statistik pengaduan untuk pengambilan keputusan.",
    ]:
        story.append(Paragraph(f"• {item}", s["bullet"]))
    story.append(Spacer(1, 0.4*cm))

    # ── 2. Deskripsi Sistem ───────────────────────────────────────────────────
    story.append(section_header("2. Deskripsi Sistem", s))
    story.append(Spacer(1, 0.3*cm))

    for label, value in [
        ("Nama Sistem",    "CAMPUSCARE — Campus Facility Complaint and Management System"),
        ("Jenis Produk",   "Backend REST API"),
        ("Bahasa",         "Go (Golang)"),
        ("Arsitektur",     "Clean Architecture"),
        ("Database",       "PostgreSQL"),
        ("Framework HTTP", "Fiber v2"),
        ("Repositori",     "https://github.com/AanFLenrir/campuscare-backend"),
    ]:
        story.append(info_box(label, value, s))
    story.append(Spacer(1, 0.5*cm))

    story.append(Paragraph(
        "CAMPUSCARE adalah sistem backend REST API berbasis Go yang dirancang dengan "
        "pendekatan Clean Architecture. Sistem ini memisahkan tanggung jawab ke dalam "
        "layer yang jelas: Model, Repository, Service, dan Route — tanpa menggunakan "
        "controller atau handler layer tambahan sesuai ketentuan tugas.",
        s["body"]
    ))
    story.append(Spacer(1, 0.4*cm))

    # ── 3. Tujuan Sistem ──────────────────────────────────────────────────────
    story.append(section_header("3. Tujuan Sistem", s))
    story.append(Spacer(1, 0.3*cm))
    for i, item in enumerate([
        "Menyediakan platform pelaporan kerusakan fasilitas yang mudah digunakan oleh mahasiswa.",
        "Mempercepat respons penanganan pengaduan oleh petugas pemeliharaan.",
        "Memberikan transparansi status penanganan kepada pelapor secara real-time.",
        "Menyediakan data statistik pengaduan untuk pengambilan keputusan manajemen.",
        "Membangun fondasi sistem yang scalable, maintainable, dan mudah dikembangkan.",
        "Menerapkan prinsip Clean Architecture dalam pengembangan backend Go.",
    ], 1):
        story.append(Paragraph(f"{i}. {item}", s["bullet"]))
    story.append(Spacer(1, 0.4*cm))

    # ── 4. Fitur Sistem ───────────────────────────────────────────────────────
    story.append(section_header("4. Fitur Sistem", s))
    story.append(Spacer(1, 0.3*cm))
    fitur_rows = [
        ["1", "Autentikasi & Otorisasi", "Register, login JWT, profil pengguna, role-based access control"],
        ["2", "Manajemen Pengguna",      "Admin dapat membuat, membaca, memperbarui, menonaktifkan pengguna"],
        ["3", "Manajemen Fasilitas",     "CRUD fasilitas kampus dengan data lokasi dan kondisi"],
        ["4", "Manajemen Kategori",      "CRUD kategori pengaduan untuk klasifikasi masalah"],
        ["5", "Pengaduan Mahasiswa",     "Mahasiswa dapat membuat dan memantau pengaduan milik sendiri"],
        ["6", "Penugasan Petugas",       "Admin dapat menugaskan petugas ke pengaduan tertentu"],
        ["7", "Pembaruan Status",        "Petugas dapat memperbarui status penanganan pengaduan"],
        ["8", "Riwayat Pengaduan",       "Rekam jejak setiap perubahan status pengaduan beserta catatannya"],
        ["9", "Laporan & Statistik",     "Statistik jumlah pengaduan berdasarkan status dan kategori"],
    ]
    story.append(make_table(
        ["No", "Fitur", "Deskripsi"],
        fitur_rows,
        col_widths=[1*cm, 4.5*cm, 10*cm]
    ))
    story.append(Spacer(1, 0.4*cm))

    # ── 5. Role Pengguna ──────────────────────────────────────────────────────
    story.append(section_header("5. Role Pengguna dan Hak Akses", s))
    story.append(Spacer(1, 0.3*cm))

    for role, hak in [
        ("ADMIN", [
            "Mengelola data pengguna (Create, Read, Update, Delete/Deactivate)",
            "Mengelola data fasilitas kampus",
            "Mengelola kategori pengaduan",
            "Melihat seluruh pengaduan dari semua mahasiswa",
            "Mengubah status pengaduan",
            "Menugaskan petugas ke pengaduan tertentu",
            "Mengakses laporan dan statistik pengaduan",
        ]),
        ("MAHASISWA", [
            "Registrasi akun baru",
            "Login dan mendapatkan token JWT",
            "Melihat profil sendiri",
            "Membuat pengaduan baru dengan prioritas dan kategori",
            "Melihat daftar pengaduan milik sendiri",
            "Melihat detail dan riwayat status pengaduan sendiri",
        ]),
        ("PETUGAS", [
            "Login dan mendapatkan token JWT",
            "Melihat daftar tugas yang diberikan admin",
            "Melihat detail pengaduan yang ditugaskan",
            "Memperbarui status penanganan pengaduan",
            "Menambahkan catatan penanganan",
        ]),
    ]:
        story.append(Paragraph(f"<b>{role}</b>", s["h3"]))
        for item in hak:
            story.append(Paragraph(f"• {item}", s["bullet"]))
        story.append(Spacer(1, 0.2*cm))
    story.append(Spacer(1, 0.3*cm))

    # ── 6. Status Pengaduan ───────────────────────────────────────────────────
    story.append(section_header("6. Status dan Prioritas Pengaduan", s))
    story.append(Spacer(1, 0.3*cm))

    story.append(Paragraph("<b>Alur Status Pengaduan:</b>", s["h3"]))
    status_rows = [
        ["OPEN",        "Pengaduan baru dibuat oleh mahasiswa"],
        ["VERIFIED",    "Pengaduan telah diverifikasi oleh admin"],
        ["ASSIGNED",    "Pengaduan telah ditugaskan ke petugas"],
        ["IN_PROGRESS", "Petugas sedang menangani pengaduan"],
        ["RESOLVED",    "Pengaduan telah diselesaikan oleh petugas"],
        ["CLOSED",      "Pengaduan ditutup setelah konfirmasi"],
        ["REJECTED",    "Pengaduan ditolak oleh admin"],
    ]
    story.append(make_table(["Status", "Deskripsi"], status_rows, col_widths=[4*cm, 11.5*cm]))
    story.append(Spacer(1, 0.3*cm))

    story.append(Paragraph("<b>Level Prioritas:</b>", s["h3"]))
    prio_rows = [
        ["LOW",    "Prioritas rendah, tidak mengganggu aktivitas"],
        ["MEDIUM", "Prioritas sedang, perlu ditangani dalam waktu dekat"],
        ["HIGH",   "Prioritas tinggi, mengganggu aktivitas kampus"],
        ["URGENT", "Sangat mendesak, perlu ditangani segera"],
    ]
    story.append(make_table(["Prioritas", "Deskripsi"], prio_rows, col_widths=[4*cm, 11.5*cm]))
    story.append(Spacer(1, 0.4*cm))
    story.append(PageBreak())

    # ── 7. Rancangan Database ─────────────────────────────────────────────────
    story.append(section_header("7. Rancangan Database", s))
    story.append(Spacer(1, 0.3*cm))
    story.append(Paragraph(
        "Database menggunakan PostgreSQL dengan 8 tabel utama. Setiap tabel menggunakan "
        "UUID sebagai primary key untuk keamanan dan skalabilitas yang lebih baik.",
        s["body"]
    ))
    story.append(Spacer(1, 0.3*cm))

    tables_data = [
        ("roles", [
            ["id",          "UUID",         "Primary Key"],
            ["name",        "VARCHAR(50)",  "ADMIN / MAHASISWA / PETUGAS"],
            ["description", "TEXT",         "Deskripsi role"],
            ["created_at",  "TIMESTAMP",    "Waktu dibuat"],
            ["updated_at",  "TIMESTAMP",    "Waktu diperbarui"],
        ]),
        ("users", [
            ["id",             "UUID",         "Primary Key"],
            ["role_id",        "UUID",         "FK → roles.id"],
            ["name",           "VARCHAR(255)", "Nama lengkap pengguna"],
            ["email",          "VARCHAR(255)", "Email unik (index)"],
            ["password",       "TEXT",         "Hash bcrypt"],
            ["phone",          "VARCHAR(20)",  "Nomor telepon"],
            ["student_number", "VARCHAR(50)",  "NIM (nullable, khusus mahasiswa)"],
            ["is_active",      "BOOLEAN",      "Status akun aktif/nonaktif"],
            ["created_at",     "TIMESTAMP",    "Waktu dibuat"],
            ["updated_at",     "TIMESTAMP",    "Waktu diperbarui"],
        ]),
        ("facilities", [
            ["id",          "UUID",         "Primary Key"],
            ["name",        "VARCHAR(255)", "Nama fasilitas"],
            ["location",    "VARCHAR(255)", "Lokasi fasilitas"],
            ["description", "TEXT",         "Deskripsi fasilitas"],
            ["condition",   "VARCHAR(50)",  "Kondisi fasilitas"],
            ["created_at",  "TIMESTAMP",    "Waktu dibuat"],
            ["updated_at",  "TIMESTAMP",    "Waktu diperbarui"],
        ]),
        ("categories", [
            ["id",          "UUID",         "Primary Key"],
            ["name",        "VARCHAR(100)", "Nama kategori pengaduan"],
            ["description", "TEXT",         "Deskripsi kategori"],
            ["created_at",  "TIMESTAMP",    "Waktu dibuat"],
            ["updated_at",  "TIMESTAMP",    "Waktu diperbarui"],
        ]),
        ("complaints", [
            ["id",          "UUID",         "Primary Key"],
            ["user_id",     "UUID",         "FK → users.id (pelapor)"],
            ["facility_id", "UUID",         "FK → facilities.id"],
            ["category_id", "UUID",         "FK → categories.id"],
            ["title",       "VARCHAR(255)", "Judul pengaduan"],
            ["description", "TEXT",         "Detail pengaduan"],
            ["priority",    "VARCHAR(20)",  "LOW / MEDIUM / HIGH / URGENT"],
            ["status",      "VARCHAR(30)",  "Status pengaduan (enum)"],
            ["created_at",  "TIMESTAMP",    "Waktu dibuat"],
            ["updated_at",  "TIMESTAMP",    "Waktu diperbarui"],
        ]),
        ("complaint_images", [
            ["id",           "UUID",      "Primary Key"],
            ["complaint_id", "UUID",      "FK → complaints.id"],
            ["image_url",    "TEXT",      "URL atau path file gambar"],
            ["created_at",   "TIMESTAMP", "Waktu diunggah"],
        ]),
        ("assignments", [
            ["id",           "UUID",      "Primary Key"],
            ["complaint_id", "UUID",      "FK → complaints.id"],
            ["officer_id",   "UUID",      "FK → users.id (petugas)"],
            ["assigned_by",  "UUID",      "FK → users.id (admin)"],
            ["assigned_at",  "TIMESTAMP", "Waktu penugasan"],
            ["notes",        "TEXT",      "Catatan penugasan (nullable)"],
        ]),
        ("complaint_histories", [
            ["id",           "UUID",      "Primary Key"],
            ["complaint_id", "UUID",      "FK → complaints.id"],
            ["user_id",      "UUID",      "FK → users.id (yang mengubah)"],
            ["old_status",   "VARCHAR(30)", "Status sebelum perubahan"],
            ["new_status",   "VARCHAR(30)", "Status setelah perubahan"],
            ["notes",        "TEXT",      "Catatan perubahan"],
            ["created_at",   "TIMESTAMP", "Waktu perubahan"],
        ]),
    ]

    for tbl_name, tbl_rows in tables_data:
        story.append(Paragraph(f"Tabel: <b>{tbl_name}</b>", s["h3"]))
        story.append(make_table(
            ["Kolom", "Tipe Data", "Keterangan"],
            tbl_rows,
            col_widths=[4.5*cm, 3.5*cm, 7.5*cm]
        ))
        story.append(Spacer(1, 0.3*cm))

    story.append(PageBreak())

    # ── 8. ERD ────────────────────────────────────────────────────────────────
    story.append(section_header("8. Entity-Relationship Diagram (ERD)", s))
    story.append(Spacer(1, 0.3*cm))

    story.append(Paragraph(
        "Diagram ERD berikut menggambarkan relasi antar entitas dalam sistem CAMPUSCARE:",
        s["body"]
    ))
    story.append(Spacer(1, 0.3*cm))

    # Draw ERD using Table as visual representation
    erd_style = ParagraphStyle(
        "erd", fontSize=9, fontName="Courier", leading=13,
        textColor=GRAY_DARK, alignment=TA_CENTER
    )

    # Visual ERD table
    erd_entities = [
        ["roles", "users", "complaints", "complaint_images"],
        ["",      "",      "assignments", ""],
        ["",      "",      "complaint_histories", ""],
        ["facilities", "", "", ""],
        ["categories", "", "", ""],
    ]

    # Build ERD as styled description
    erd_rows = [
        ["roles",               "1 : N", "users",               "Setiap role dapat dimiliki banyak pengguna"],
        ["users",               "1 : N", "complaints",          "Satu mahasiswa dapat memiliki banyak pengaduan"],
        ["users",               "1 : N", "assignments",         "Satu petugas dapat menerima banyak penugasan"],
        ["users",               "1 : N", "complaint_histories", "Satu pengguna dapat membuat banyak riwayat"],
        ["facilities",          "1 : N", "complaints",          "Satu fasilitas dapat memiliki banyak pengaduan"],
        ["categories",          "1 : N", "complaints",          "Satu kategori dapat memiliki banyak pengaduan"],
        ["complaints",          "1 : N", "complaint_images",    "Satu pengaduan dapat memiliki banyak gambar"],
        ["complaints",          "1 : N", "assignments",         "Satu pengaduan dapat ditugaskan ke banyak petugas"],
        ["complaints",          "1 : N", "complaint_histories", "Satu pengaduan memiliki banyak riwayat perubahan"],
    ]
    story.append(make_table(
        ["Entitas Asal", "Relasi", "Entitas Tujuan", "Keterangan"],
        erd_rows,
        col_widths=[3.5*cm, 2*cm, 4*cm, 6*cm]
    ))
    story.append(Spacer(1, 0.4*cm))

    # ERD text diagram
    story.append(Paragraph("<b>Diagram Relasi (Notasi Teks):</b>", s["h3"]))
    erd_diagram = [
        [Paragraph(
            "roles ──────────────── users\n"
            "                          │\n"
            "                    ┌─────┼──────────────────────────────┐\n"
            "                    │     │                              │\n"
            "             complaints  complaint_histories       assignments\n"
            "                    │\n"
            "          ┌─────────┼──────────┐\n"
            "          │         │          │\n"
            "   complaint_   assignments  complaint_\n"
            "    images                  histories\n\n"
            "facilities ──── complaints\n"
            "categories ──── complaints",
            ParagraphStyle(
                "erd_diag", fontSize=8.5, fontName="Courier", leading=13,
                textColor=GRAY_DARK, alignment=TA_LEFT
            )
        )]
    ]
    erd_tbl = Table(erd_diagram, colWidths=[A4[0] - 5*cm])
    erd_tbl.setStyle(TableStyle([
        ("BACKGROUND", (0, 0), (-1, -1), GRAY_LIGHT),
        ("LINEABOVE",  (0, 0), (-1, -1), 1, BLUE_DARK),
        ("LINEBELOW",  (0, 0), (-1, -1), 1, BLUE_DARK),
        ("LINEBEFORE", (0, 0), (-1, -1), 1, BLUE_DARK),
        ("LINEAFTER",  (0, 0), (-1, -1), 1, BLUE_DARK),
        ("TOPPADDING", (0, 0), (-1, -1), 12),
        ("BOTTOMPADDING", (0, 0), (-1, -1), 12),
        ("LEFTPADDING", (0, 0), (-1, -1), 15),
        ("RIGHTPADDING", (0, 0), (-1, -1), 15),
    ]))
    story.append(erd_tbl)
    story.append(Spacer(1, 0.4*cm))
    story.append(PageBreak())

    # ── 9. Arsitektur Clean Architecture ──────────────────────────────────────
    story.append(section_header("9. Arsitektur Clean Architecture", s))
    story.append(Spacer(1, 0.3*cm))

    story.append(Paragraph(
        "Sistem CAMPUSCARE dibangun menggunakan pendekatan Clean Architecture yang memisahkan "
        "tanggung jawab ke dalam layer yang independen satu sama lain. Ketergantungan "
        "berjalan dari layer luar (HTTP) ke dalam (domain), tidak sebaliknya.",
        s["body"]
    ))
    story.append(Spacer(1, 0.3*cm))

    # Architecture layers visual
    arch_layers = [
        ("HTTP Layer  ·  route/", "Menerima HTTP request, memanggil service, mengirim response.\nMenggunakan Fiber v2. Tidak mengandung business logic.", BLUE_DARK),
        ("Service Layer  ·  app/service/", "Business logic murni. Tidak bergantung pada Fiber.\nMenerima context.Context dan parameter bertipe jelas.\nMengembalikan domain object dan error.", BLUE_MID),
        ("Repository Layer  ·  app/repository/", "Akses database PostgreSQL menggunakan pgx.\nImplementasi interface dengan parameterized query.\nTidak mengandung business logic.", ACCENT),
        ("Model Layer  ·  app/model/", "Struct Go yang merepresentasikan entitas domain.\nTidak mengandung business logic atau dependensi eksternal.", colors.HexColor("#5c6bc0")),
    ]

    for label, desc, color in arch_layers:
        layer_data = [[
            Paragraph(f"<b>{label}</b>", ParagraphStyle(
                "al", fontSize=10, textColor=WHITE, fontName="Helvetica-Bold", leading=14
            )),
            Paragraph(desc, ParagraphStyle(
                "ad", fontSize=9.5, textColor=GRAY_DARK, fontName="Helvetica", leading=14
            ))
        ]]
        lt = Table(layer_data, colWidths=[5*cm, 10.5*cm])
        lt.setStyle(TableStyle([
            ("BACKGROUND", (0, 0), (0, 0), color),
            ("BACKGROUND", (1, 0), (1, 0), BLUE_LIGHT),
            ("VALIGN", (0, 0), (-1, -1), "MIDDLE"),
            ("TOPPADDING", (0, 0), (-1, -1), 8),
            ("BOTTOMPADDING", (0, 0), (-1, -1), 8),
            ("LEFTPADDING", (0, 0), (-1, -1), 10),
            ("LINEBELOW", (0, 0), (-1, -1), 2, WHITE),
        ]))
        story.append(lt)

    story.append(Spacer(1, 0.3*cm))

    # Arrow indicator
    arrow_data = [["▼  Ketergantungan mengalir dari HTTP Layer ke Model Layer  ▼"]]
    at = Table(arrow_data, colWidths=[A4[0] - 5*cm])
    at.setStyle(TableStyle([
        ("BACKGROUND", (0, 0), (-1, -1), colors.HexColor("#e3f2fd")),
        ("FONTNAME",   (0, 0), (-1, -1), "Helvetica-Bold"),
        ("FONTSIZE",   (0, 0), (-1, -1), 9),
        ("ALIGN",      (0, 0), (-1, -1), "CENTER"),
        ("TEXTCOLOR",  (0, 0), (-1, -1), BLUE_DARK),
        ("TOPPADDING", (0, 0), (-1, -1), 6),
        ("BOTTOMPADDING", (0, 0), (-1, -1), 6),
    ]))
    story.append(at)
    story.append(Spacer(1, 0.4*cm))

    story.append(Paragraph("<b>Prinsip yang Diterapkan:</b>", s["h3"]))
    for item in [
        "Service tidak menerima *fiber.Ctx — HTTP parsing dilakukan di layer route.",
        "Repository hanya berinteraksi dengan database melalui interface.",
        "Model hanya berisi struct Go tanpa business logic.",
        "Middleware JWT dan Authorization berada di layer route/middleware.",
        "Dependency injection dilakukan di route.Setup() saat aplikasi dimulai.",
    ]:
        story.append(Paragraph(f"• {item}", s["bullet"]))
    story.append(Spacer(1, 0.4*cm))

    story.append(Paragraph("<b>Komponen Pendukung:</b>", s["h3"]))
    comp_rows = [
        ["config/",     "Konfigurasi aplikasi dari environment (.env)"],
        ["database/",   "Koneksi pool PostgreSQL menggunakan pgx"],
        ["helper/",     "Fungsi JWT, bcrypt, dan HTTP response formatter"],
        ["middleware/",  "JWT authentication dan role-based authorization"],
        ["migrations/",  "File SQL migrasi database (up dan down)"],
        ["docs/",        "Dokumentasi sistem, perancangan, dan laporan"],
    ]
    story.append(make_table(
        ["Direktori", "Fungsi"],
        comp_rows,
        col_widths=[4*cm, 11.5*cm]
    ))
    story.append(Spacer(1, 0.4*cm))
    story.append(PageBreak())

    # ── 10. Teknologi ─────────────────────────────────────────────────────────
    story.append(section_header("10. Teknologi yang Digunakan", s))
    story.append(Spacer(1, 0.3*cm))
    tech_rows = [
        ["Go (Golang)",             "1.21+",    "Bahasa pemrograman utama backend"],
        ["Fiber v2",                "v2.52+",   "HTTP framework berbasis fasthttp"],
        ["PostgreSQL",              "14+",      "Sistem manajemen basis data relasional"],
        ["pgx v5",                  "v5",       "Driver PostgreSQL native untuk Go"],
        ["golang-jwt",              "v5",       "Library JWT untuk autentikasi"],
        ["golang.org/x/crypto",     "latest",   "Bcrypt untuk hashing password"],
        ["google/uuid",             "v1",       "Pembuatan UUID v4"],
        ["gofmt",                   "bawaan",   "Format kode Go secara otomatis"],
    ]
    story.append(make_table(
        ["Teknologi", "Versi", "Fungsi"],
        tech_rows,
        col_widths=[5*cm, 2.5*cm, 8*cm]
    ))
    story.append(Spacer(1, 0.4*cm))

    # ── 11. Struktur Project ──────────────────────────────────────────────────
    story.append(section_header("11. Struktur Project", s))
    story.append(Spacer(1, 0.3*cm))

    struct_text = (
        "campuscare-backend/\n"
        "├── app/\n"
        "│   ├── model/          # Entitas domain (8 struct)\n"
        "│   ├── repository/     # Akses database (7 repository)\n"
        "│   └── service/        # Business logic (7 service + 1 test file)\n"
        "├── config/             # Konfigurasi dari .env\n"
        "├── database/           # Koneksi PostgreSQL pool\n"
        "├── docs/\n"
        "│   ├── perancangan/    # Dokumen rancangan sistem\n"
        "│   ├── report/         # Laporan akhir (file ini)\n"
        "│   ├── git-workflow.md\n"
        "│   └── team-contribution.md\n"
        "├── helper/             # JWT, bcrypt, response formatter\n"
        "├── middleware/         # Auth & authorization middleware\n"
        "├── migrations/         # SQL migration (8 tabel, up+down)\n"
        "├── route/              # HTTP endpoint registration\n"
        "├── .env.example        # Template konfigurasi\n"
        "├── .gitignore\n"
        "├── go.mod\n"
        "├── go.sum\n"
        "├── main.go             # Entry point aplikasi\n"
        "└── README.md"
    )
    struct_data = [[Paragraph(struct_text, ParagraphStyle(
        "struct", fontSize=8.5, fontName="Courier", leading=13,
        textColor=GRAY_DARK, alignment=TA_LEFT
    ))]]
    st = Table(struct_data, colWidths=[A4[0] - 5*cm])
    st.setStyle(TableStyle([
        ("BACKGROUND", (0, 0), (-1, -1), GRAY_LIGHT),
        ("LINEABOVE",  (0, 0), (-1, -1), 1.5, BLUE_DARK),
        ("LINEBELOW",  (0, 0), (-1, -1), 1.5, BLUE_DARK),
        ("LINEBEFORE", (0, 0), (-1, -1), 1.5, BLUE_DARK),
        ("LINEAFTER",  (0, 0), (-1, -1), 1.5, BLUE_DARK),
        ("TOPPADDING", (0, 0), (-1, -1), 12),
        ("BOTTOMPADDING", (0, 0), (-1, -1), 12),
        ("LEFTPADDING", (0, 0), (-1, -1), 15),
    ]))
    story.append(st)
    story.append(Spacer(1, 0.4*cm))
    story.append(PageBreak())

    # ── 12. Rancangan Endpoint API ────────────────────────────────────────────
    story.append(section_header("12. Rancangan Endpoint API", s))
    story.append(Spacer(1, 0.3*cm))

    story.append(Paragraph("Base URL: <b>http://localhost:3000/api/v1</b>", s["body_left"]))
    story.append(Spacer(1, 0.3*cm))

    endpoint_groups = [
        ("12.1 Authentication (Public)", [
            ["POST", "/auth/register",   "—",          "Register mahasiswa baru"],
            ["POST", "/auth/login",      "—",          "Login dan dapatkan JWT token"],
            ["GET",  "/auth/me",         "JWT",        "Lihat profil pengguna"],
            ["GET",  "/health",          "—",          "Health check server"],
        ]),
        ("12.2 Manajemen Pengguna (ADMIN)", [
            ["GET",    "/users",     "ADMIN + JWT", "Daftar semua pengguna"],
            ["GET",    "/users/:id", "ADMIN + JWT", "Detail pengguna berdasarkan ID"],
            ["POST",   "/users",     "ADMIN + JWT", "Buat pengguna baru"],
            ["PUT",    "/users/:id", "ADMIN + JWT", "Perbarui data pengguna"],
            ["DELETE", "/users/:id", "ADMIN + JWT", "Nonaktifkan pengguna"],
        ]),
        ("12.3 Fasilitas", [
            ["GET",    "/facilities",     "JWT",        "Daftar semua fasilitas"],
            ["GET",    "/facilities/:id", "JWT",        "Detail fasilitas"],
            ["POST",   "/facilities",     "ADMIN + JWT", "Tambah fasilitas baru"],
            ["PUT",    "/facilities/:id", "ADMIN + JWT", "Perbarui fasilitas"],
            ["DELETE", "/facilities/:id", "ADMIN + JWT", "Hapus fasilitas"],
        ]),
        ("12.4 Kategori", [
            ["GET",    "/categories",     "JWT",        "Daftar semua kategori"],
            ["GET",    "/categories/:id", "JWT",        "Detail kategori"],
            ["POST",   "/categories",     "ADMIN + JWT", "Tambah kategori baru"],
            ["PUT",    "/categories/:id", "ADMIN + JWT", "Perbarui kategori"],
            ["DELETE", "/categories/:id", "ADMIN + JWT", "Hapus kategori"],
        ]),
        ("12.5 Pengaduan — Mahasiswa", [
            ["POST", "/my-complaints",     "MAHASISWA + JWT", "Buat pengaduan baru"],
            ["GET",  "/my-complaints",     "MAHASISWA + JWT", "Daftar pengaduan sendiri"],
            ["GET",  "/my-complaints/:id", "MAHASISWA + JWT", "Detail pengaduan sendiri"],
        ]),
        ("12.6 Pengaduan — Admin", [
            ["GET",  "/complaints",           "ADMIN + JWT", "Semua pengaduan"],
            ["GET",  "/complaints/:id",        "ADMIN + JWT", "Detail pengaduan"],
            ["PUT",  "/complaints/:id/status", "ADMIN + JWT", "Ubah status pengaduan"],
            ["POST", "/complaints/:id/assign", "ADMIN + JWT", "Tugaskan petugas"],
        ]),
        ("12.7 Penugasan — Petugas", [
            ["GET", "/assignments",           "PETUGAS + JWT", "Daftar tugas yang diterima"],
            ["PUT", "/assignments/:id/status", "PETUGAS + JWT", "Perbarui status pengaduan"],
        ]),
        ("12.8 Laporan — Admin", [
            ["GET", "/reports/summary", "ADMIN + JWT", "Statistik ringkasan pengaduan"],
        ]),
    ]

    for group_name, rows in endpoint_groups:
        story.append(Paragraph(group_name, s["h3"]))
        story.append(make_table(
            ["Method", "Endpoint", "Auth", "Deskripsi"],
            rows,
            col_widths=[1.8*cm, 5.5*cm, 3.2*cm, 5*cm]
        ))
        story.append(Spacer(1, 0.3*cm))

    story.append(Spacer(1, 0.2*cm))
    story.append(Paragraph("<b>Format Response API:</b>", s["h3"]))
    resp_text = (
        "# Success:\n"
        '{ "success": true, "data": { ... } }\n\n'
        "# Error:\n"
        '{ "success": false, "message": "deskripsi error" }'
    )
    resp_data = [[Paragraph(resp_text, ParagraphStyle(
        "resp", fontSize=9, fontName="Courier", leading=13, textColor=GRAY_DARK
    ))]]
    rt = Table(resp_data, colWidths=[A4[0] - 5*cm])
    rt.setStyle(TableStyle([
        ("BACKGROUND", (0, 0), (-1, -1), GRAY_LIGHT),
        ("LINEABOVE",  (0, 0), (-1, -1), 1, BLUE_DARK),
        ("LINEBELOW",  (0, 0), (-1, -1), 1, BLUE_DARK),
        ("LINEBEFORE", (0, 0), (-1, -1), 1, BLUE_DARK),
        ("LINEAFTER",  (0, 0), (-1, -1), 1, BLUE_DARK),
        ("TOPPADDING", (0, 0), (-1, -1), 10),
        ("BOTTOMPADDING", (0, 0), (-1, -1), 10),
        ("LEFTPADDING", (0, 0), (-1, -1), 15),
    ]))
    story.append(rt)
    story.append(Spacer(1, 0.4*cm))
    story.append(PageBreak())

    # ── 13. Implementasi Boilerplate ──────────────────────────────────────────
    story.append(section_header("13. Implementasi Boilerplate", s))
    story.append(Spacer(1, 0.3*cm))

    impl_rows = [
        ["Model Layer",      "8 entitas: Role, User, Facility, Category, Complaint, ComplaintImage, Assignment, ComplaintHistory"],
        ["Repository Layer", "7 repository dengan interface dan implementasi PostgreSQL menggunakan parameterized query"],
        ["Service Layer",    "7 service dengan pure business logic (tanpa *fiber.Ctx):\nAuthService, UserService, FacilityService, CategoryService,\nComplaintService, AssignmentService, ReportService"],
        ["Route Layer",      "HTTP callbacks inline di route/route.go yang menghubungkan\nFiber request ke service layer"],
        ["Middleware",       "JWT Authentication (authentication.go) dan\nRole-based Authorization (authorization.go)"],
        ["Helper",           "helper/jwt.go: GenerateToken, ValidateToken\nhelper/password.go: HashPassword, CheckPassword\nhelper/response.go: Success, Created, Error"],
        ["Migration",        "8 file migrasi UP + 8 file migrasi DOWN\n(000001 s/d 000008 untuk setiap tabel)"],
    ]
    for label, desc in impl_rows:
        story.append(Paragraph(f"<b>{label}</b>", s["h3"]))
        story.append(Paragraph(desc, s["body_left"]))
        story.append(Spacer(1, 0.2*cm))

    story.append(Spacer(1, 0.4*cm))

    # ── 14. Panduan Instalasi ─────────────────────────────────────────────────
    story.append(section_header("14. Panduan Instalasi dan Menjalankan", s))
    story.append(Spacer(1, 0.3*cm))

    story.append(Paragraph("<b>Prasyarat:</b>", s["h3"]))
    for item in ["Go 1.21 atau lebih baru", "PostgreSQL 14+", "Git"]:
        story.append(Paragraph(f"• {item}", s["bullet"]))
    story.append(Spacer(1, 0.3*cm))

    setup_steps = [
        ("1. Clone Repository",
         "git clone https://github.com/AanFLenrir/campuscare-backend.git\ncd campuscare-backend"),
        ("2. Konfigurasi Environment",
         "cp .env.example .env\n# Edit .env sesuai konfigurasi database dan JWT secret Anda"),
        ("3. Download Dependencies",
         "go mod tidy"),
        ("4. Jalankan Migrasi Database",
         "# Jalankan file SQL di folder migrations/ secara berurutan\npsql -U postgres -d campuscare -f migrations/000001_create_roles_table.up.sql\n# dst. hingga 000008_..."),
        ("5. Jalankan Aplikasi",
         "go run main.go\n# Aplikasi berjalan di http://localhost:3000"),
        ("6. Verifikasi",
         "curl http://localhost:3000/api/v1/health\n# Response: {\"status\":\"ok\"}"),
    ]

    for title, cmd in setup_steps:
        story.append(Paragraph(title, s["h3"]))
        cmd_data = [[Paragraph(cmd, ParagraphStyle(
            "cmd", fontSize=9, fontName="Courier", leading=13, textColor=GRAY_DARK
        ))]]
        ct = Table(cmd_data, colWidths=[A4[0] - 5*cm])
        ct.setStyle(TableStyle([
            ("BACKGROUND", (0, 0), (-1, -1), GRAY_LIGHT),
            ("LINEBEFORE", (0, 0), (-1, -1), 3, BLUE_DARK),
            ("TOPPADDING", (0, 0), (-1, -1), 8),
            ("BOTTOMPADDING", (0, 0), (-1, -1), 8),
            ("LEFTPADDING", (0, 0), (-1, -1), 12),
        ]))
        story.append(ct)
        story.append(Spacer(1, 0.2*cm))

    story.append(Spacer(1, 0.4*cm))
    story.append(PageBreak())

    # ── 15. Pengujian ─────────────────────────────────────────────────────────
    story.append(section_header("15. Pengujian", s))
    story.append(Spacer(1, 0.3*cm))

    story.append(Paragraph("<b>15.1 Unit Test</b>", s["h2"]))
    story.append(Paragraph(
        "Unit test ditulis untuk ComplaintService menggunakan mock repository "
        "tanpa memerlukan koneksi database. File: app/service/complaint_service_test.go",
        s["body"]
    ))
    story.append(Spacer(1, 0.2*cm))

    test_rows = [
        ["TestCreateComplaint_Success",       "Pembuatan pengaduan dengan input valid",          "✓ PASS"],
        ["TestCreateComplaint_MissingFields",  "Validasi input wajib yang tidak lengkap",         "✓ PASS"],
        ["TestListMyComplaints",               "Daftar pengaduan milik pengguna tertentu",        "✓ PASS"],
        ["TestGetMyComplaint_Forbidden",       "Akses pengaduan milik pengguna lain (forbidden)", "✓ PASS"],
        ["TestUpdateComplaintStatus",          "Pembaruan status pengaduan oleh petugas",         "✓ PASS"],
        ["TestAssignComplaint",                "Penugasan petugas ke pengaduan oleh admin",        "✓ PASS"],
    ]
    story.append(make_table(
        ["Nama Test", "Skenario", "Hasil"],
        test_rows,
        col_widths=[5.5*cm, 7*cm, 2*cm]
    ))
    story.append(Spacer(1, 0.3*cm))

    story.append(Paragraph("<b>15.2 Build dan Static Analysis</b>", s["h2"]))
    build_rows = [
        ["gofmt -w .",    "Format kode otomatis",           "✓ OK — tidak ada perubahan"],
        ["go mod tidy",   "Bersihkan dependency",           "✓ OK"],
        ["go vet ./...",  "Static analysis",                "✓ OK — tidak ada isu"],
        ["go test ./...", "Jalankan semua unit test",       "✓ OK — 6/6 PASS"],
        ["go build ./...", "Kompilasi seluruh package",     "✓ OK — binary dikompilasi"],
    ]
    story.append(make_table(
        ["Perintah", "Deskripsi", "Hasil"],
        build_rows,
        col_widths=[4*cm, 5.5*cm, 6*cm]
    ))
    story.append(Spacer(1, 0.3*cm))

    story.append(Paragraph("<b>15.3 Integration Test</b>", s["h2"]))
    story.append(Paragraph(
        "Integration test memerlukan koneksi PostgreSQL aktif. "
        "Jalankan setelah konfigurasi database dengan perintah: go test ./... -tags integration",
        s["body"]
    ))
    story.append(Spacer(1, 0.4*cm))

    # ── 16. GitHub ────────────────────────────────────────────────────────────
    story.append(section_header("16. Repository GitHub dan Git Workflow", s))
    story.append(Spacer(1, 0.3*cm))

    for label, value in [
        ("URL Repository", "https://github.com/AanFLenrir/campuscare-backend"),
        ("Branch Utama",   "main"),
        ("Visibilitas",    "Public (harus dikonfirmasi di GitHub Settings)"),
        ("Total Commit",   "4 commit (first commit + 3 commit implementasi)"),
        ("Git Master",     "MOHAMMAD ANHAZ ABDILAH (AanFLenrir)"),
    ]:
        story.append(info_box(label, value, s))
    story.append(Spacer(1, 0.3*cm))

    story.append(Paragraph("<b>Riwayat Commit:</b>", s["h3"]))
    commit_rows = [
        ["796a90f", "chore: allow .env.example to be tracked, update .gitignore",                     "MOHAMMAD ANHAZ ABDILAH"],
        ["b516d9d", "docs: update git-workflow, team-contribution, project-report",                   "MOHAMMAD ANHAZ ABDILAH"],
        ["0695298", "feat: implement full CAMPUSCARE backend with clean architecture",                 "MOHAMMAD ANHAZ ABDILAH"],
        ["67c22dd", "first commit",                                                                    "MOHAMMAD ANHAZ ABDILAH"],
    ]
    story.append(make_table(
        ["Hash", "Pesan Commit", "Author"],
        commit_rows,
        col_widths=[2*cm, 9.5*cm, 4*cm]
    ))
    story.append(Spacer(1, 0.3*cm))

    story.append(Paragraph("<b>Konvensi Micro Commits:</b>", s["h3"]))
    for item in [
        "Setiap commit fokus pada satu perubahan yang dapat dijelaskan dalam satu kalimat.",
        "Format: <type>(<scope>): <deskripsi> (feat, fix, refactor, docs, test, chore)",
        "Setiap anggota commit menggunakan akun GitHub masing-masing.",
        "Git Master (MOHAMMAD ANHAZ ABDILAH) bertanggung jawab atas merging ke main.",
        "Dilarang menggunakan force push ke branch main.",
        "File .env dan credential tidak boleh dimasukkan ke repository.",
    ]:
        story.append(Paragraph(f"• {item}", s["bullet"]))
    story.append(Spacer(1, 0.4*cm))
    story.append(PageBreak())

    # ── 17. Catatan Penting ───────────────────────────────────────────────────
    story.append(section_header("17. Catatan Penting dan Keamanan", s))
    story.append(Spacer(1, 0.3*cm))

    story.append(Paragraph("<b>Keamanan yang Diterapkan:</b>", s["h3"]))
    for item in [
        "Password di-hash menggunakan bcrypt sebelum disimpan ke database.",
        "Autentikasi menggunakan JWT (JSON Web Token) dengan expiry time.",
        "Role-based Access Control (RBAC) untuk setiap endpoint.",
        "Parameterized query untuk mencegah SQL Injection.",
        "File .env tidak dimasukkan ke repository GitHub.",
        "UUID digunakan sebagai primary key untuk mencegah enumeration attack.",
    ]:
        story.append(Paragraph(f"• {item}", s["bullet"]))
    story.append(Spacer(1, 0.3*cm))

    story.append(Paragraph("<b>Keterbatasan dan Pengembangan Lanjutan:</b>", s["h3"]))
    for item in [
        "Upload gambar pengaduan belum terimplementasi (memerlukan storage service eksternal seperti S3 atau local disk).",
        "Integration test memerlukan koneksi PostgreSQL aktif.",
        "Rate limiting belum diimplementasikan.",
        "Logging middleware belum ditambahkan.",
        "Swagger/OpenAPI documentation belum dibuat.",
        "Identitas anggota kelompok 2–4 harus dilengkapi dengan data asli.",
    ]:
        story.append(Paragraph(f"• {item}", s["bullet"]))
    story.append(Spacer(1, 0.4*cm))

    # ── Penutup ───────────────────────────────────────────────────────────────
    story.append(HRFlowable(width="100%", thickness=1, color=BLUE_DARK))
    story.append(Spacer(1, 0.4*cm))
    story.append(Paragraph(
        "Laporan ini dibuat berdasarkan implementasi aktual pada repository CAMPUSCARE Backend. "
        "Seluruh kode yang didokumentasikan telah melalui verifikasi build, static analysis, "
        "dan unit testing.",
        s["body"]
    ))
    story.append(Spacer(1, 0.3*cm))
    story.append(Paragraph(
        "Link Repository: https://github.com/AanFLenrir/campuscare-backend",
        ParagraphStyle(
            "link", fontSize=10, textColor=ACCENT, fontName="Helvetica-Bold",
            leading=14, alignment=TA_CENTER
        )
    ))
    story.append(Spacer(1, 0.5*cm))

    # Signature
    sig_data = [[
        Paragraph(
            "MOHAMMAD ANHAZ ABDILAH<br/>NIM: 434241102<br/>D4 Teknik Informatika — Universitas Airlangga",
            ParagraphStyle(
                "sig", fontSize=10, fontName="Helvetica-Bold",
                textColor=GRAY_DARK, leading=15, alignment=TA_CENTER
            )
        )
    ]]
    sig_tbl = Table(sig_data, colWidths=[A4[0] - 5*cm])
    sig_tbl.setStyle(TableStyle([
        ("BACKGROUND", (0, 0), (-1, -1), BLUE_LIGHT),
        ("LINEABOVE", (0, 0), (-1, -1), 2, BLUE_DARK),
        ("TOPPADDING", (0, 0), (-1, -1), 12),
        ("BOTTOMPADDING", (0, 0), (-1, -1), 12),
    ]))
    story.append(sig_tbl)

    return story

# ── Main ──────────────────────────────────────────────────────────────────────
if __name__ == "__main__":
    print(f"Generating PDF: {OUTPUT_PATH}")
    doc = make_doc(OUTPUT_PATH)
    styles = build_styles()
    story = build_story(styles)
    doc.build(story, onFirstPage=on_page, onLaterPages=on_page)
    size = os.path.getsize(OUTPUT_PATH)
    print(f"✓ PDF berhasil dibuat: {OUTPUT_PATH}")
    print(f"  Ukuran: {size:,} bytes ({size/1024:.1f} KB)")

