# 1. Record Architecture Decisions

## Status
Accepted

## Context
Perlu cara terdokumentasi buat mencatat keputusan desain penting (kenapa pilih X daripada Y), biar gak hilang di history chat/kepala doang.

## Decision
Pakai format ADR (Architecture Decision Record) sederhana untuk tiap keputusan signifikan — 1 file per keputusan, di `docs/adr/`.

## Consequences
Setiap keputusan besar (pilih database, pola arsitektur, message broker, dll) harus punya ADR sebelum diimplementasi.
