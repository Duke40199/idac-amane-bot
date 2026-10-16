package repository

// pgMaxBulkRows caps how many rows a single bulk INSERT sends. Postgres allows at
// most 65535 bind parameters per statement; 1000 rows stays well under that even
// for a wide row while keeping round-trips low.
const pgMaxBulkRows = 1000
