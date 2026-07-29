# File Format Coverage

This folder aims to carry one sample file per document/spreadsheet/presentation/
image/archive format that oCIS, Collabora, and ONLYOFFICE encounter in the wild.
This note tracks what's covered and what's deliberately left out.

## Covered

Word processing: doc**x**, doc**m**, dot**x**, dot**m**, docxf, odt, ott, oth,
fodt, rtf
Spreadsheets: xlsx, xlsm, xltx, xltm, ods, ots, fods, csv, tsv, dif, slk, dbf,
gnumeric, gsheet
Presentations: pptx, pptm, potx, potm, ppsx, ppsm, odp, otp, fodp, gslides
Drawings/CAD: odg, otg, fodg, svg, dxf
PDF: pdf, a PDF/A-flavored variant (`Sample-Archival-PDFA.pdf`)
Text/markup: txt, md, html, htm, xhtml, mhtml, xml, json
E-books: epub, fb2
Extensions: oxt (LibreOffice extension package)
Legacy binary (hand-built, byte-for-byte): dbf, pdb, emf, wmf
Images/video/other: jpg, jpeg, png, gif, bmp, tiff, webp, ico, url, avi, mkv,
mov, mp4, webm

## Deliberately skipped

The formats below need their originating proprietary application (or a
converter we don't have available, e.g. LibreOffice/`soffice` headless
conversion) to produce a genuinely valid file. Rather than ship a file that
merely has the right extension but invalid/garbage internal structure, they're
left out:

- **Legacy OLE binary Microsoft Office** (pre-2007): doc, dot, xls, xlsb, xlt,
  ppt, pot, pps, pub, wri
- **CorelDRAW**: cdr
- **WordPerfect**: wpd
- **Hangul Word Processor**: hwp, hwp2
- **Adobe PageMaker**: p65
- **Lotus Ami Pro**: mw
- **Sony BroadBand eBook (Librie/LRF)**: lrf
- **DjVu**: djvu
- **Adobe Flash**: swf
- **Apple iWork** (proprietary zip/IWA bundle): key, numbers
- **CGM (Computer Graphics Metafile)**: cgm
- **Microsoft Visio** (legacy binary and modern OOXML-based): vsd, vsdm, vsdx,
  vss, vssm, vssx, vstm, vstx
- **Microsoft XPS / OpenXPS**: xps, oxps
- **ONLYOFFICE form (built from docxf)**: oform
- **Kingsoft/WPS Office legacy binary**: et, ett, dps, dpt, wps, wpt
- **LibreOffice Base**: odb
- **OpenDocument chart / master document (and its template)**: odc, odm, otm
- **StarOffice 5/6 legacy XML family**: sxw, sxc, sxi, sxd, sxg, stw, stc,
  sti, std
- **HEIF/HEIC**: heif (no HEIF encoder available in this environment)

Generated 2026-07-29 for the oCIS format showcase.
