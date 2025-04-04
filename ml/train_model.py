import sqlite3
import pickle
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.neighbors import NearestNeighbors

# Подключаемся к SQLite
conn = sqlite3.connect("countries.db")
cursor = conn.cursor()

# Загружаем common_name
cursor.execute("SELECT DISTINCT common_name FROM countries")
rows = cursor.fetchall()
country_names = [row[0] for row in rows if row[0]]

# Обучаем модель
vectorizer = TfidfVectorizer(analyzer='char_wb', ngram_range=(2, 4))
X = vectorizer.fit_transform(country_names)
model = NearestNeighbors(n_neighbors=5, metric='cosine').fit(X)

# Сохраняем модель
with open("ml/country_model.pkl", "wb") as f:
    pickle.dump((model, vectorizer, country_names), f)

print("✅ Training model - Completed!")
