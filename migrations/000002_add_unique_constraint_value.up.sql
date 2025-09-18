DELETE FROM
    records
WHERE
    id NOT IN (
        SELECT
            MIN(id)
        FROM
            records
        GROUP BY
            value
    );

ALTER TABLE
    records
ADD
    CONSTRAINT unique_records_value UNIQUE (value);
