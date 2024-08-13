import sqlite3


def get_port():
    conn = sqlite3.connect("/config/db/database.db")
    c = conn.cursor()
    c.execute("SELECT value FROM settings WHERE id = 'port'")
    port = c.fetchone()[0]
    return port