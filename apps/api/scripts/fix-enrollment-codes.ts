// Script to fix enrollment codes for existing courses
import { Pool } from 'pg';

const pool = new Pool({
    host: process.env.POSTGRES_HOST || 'localhost',
    port: Number(process.env.POSTGRES_PORT || 5432),
    user: process.env.POSTGRES_USER || 'juez',
    password: process.env.POSTGRES_PASSWORD || 'juez',
    database: process.env.POSTGRES_DB || 'juez_db',
});

async function fixEnrollmentCodes() {
    try {
        console.log('Fixing enrollment codes for existing courses...');

        const result = await pool.query(`
      UPDATE courses
      SET enrollment_code = CONCAT(
        UPPER(code),
        '-',
        REPLACE(period, '-', ''),
        'G',
        group_number
      )
      WHERE enrollment_code IS NOT NULL
      RETURNING id, name, code, period, group_number, enrollment_code
    `);

        console.log(`✅ Updated ${result.rowCount} course(s):`);
        result.rows.forEach(row => {
            console.log(`  - ${row.name} (${row.code}): ${row.enrollment_code}`);
        });

    } catch (error) {
        console.error('❌ Error fixing enrollment codes:', error);
    } finally {
        await pool.end();
    }
}

fixEnrollmentCodes();
