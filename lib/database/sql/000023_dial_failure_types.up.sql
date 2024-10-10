INSERT INTO dial_error_types (id, name)
SELECT v.id, v.name
FROM (VALUES
    (0, 'Unknown'),
    (1, 'I/O Timeout'),
    (2, 'Connection Refused'),
    (3, 'Relay Circuit Failed'),
    (4, 'Relay No Reservation'),
    (5, 'Security Negotiation Failed'),
    (6, 'Concurrent Dial Succeeded'),
    (7, 'Concurrent Dial Failed'),
    (8, 'Connections Per IP Limit Exceeded'),
    (9, 'Stream Reset'),
    (10, 'Relay Resource Limit Exceeded'),
    (11, 'Error Opening Hop Stream to Relay'),
    (12, 'Dial Backoff')
) AS v(id, name)
WHERE NOT EXISTS (
    SELECT 1 FROM dial_error_types WHERE id = v.id
);