.PHONY: ml-train ml-test clean

VENV := source ml/venv/bin/activate

ml-train:
	@echo "🚀 Learn ML..."
	$(VENV) && python3 ml/train_model.py

ml-test:
	@echo "🔍 Test ML-suggestion:"
	$(VENV) && python3 ml/ml_suggest.py "georiga"

clean:
	@echo "🧹 Remove model..."
	rm -f ml/country_model.pkl
