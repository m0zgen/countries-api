import sys
import pickle
import json

with open("ml/country_model.pkl", "rb") as f:
    model, vectorizer, country_names = pickle.load(f)

def suggest_country(query):
    q_vec = vectorizer.transform([query])
    distances, indices = model.kneighbors(q_vec)
    suggestions = [country_names[i] for i in indices[0]]
    return suggestions

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print(json.dumps([]))
        sys.exit(0)

    query = sys.argv[1]
    print(json.dumps(suggest_country(query)))
