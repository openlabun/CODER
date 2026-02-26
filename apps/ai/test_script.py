"""
test_script.py
Envía el prompt de apps/ai/prompts/exercise_query.txt al LLM que esté
corriendo en Ollama (localhost:11434) y guarda los resultados en
<nombre_de_modelo>.json dentro del mismo directorio.
"""

import json
import time
import sys
from pathlib import Path

try:
    import requests
except ImportError:
    sys.exit("Instala la dependencia: pip install requests")

# ---------------------------------------------------------------------------
# Configuración
# ---------------------------------------------------------------------------
OLLAMA_BASE = "http://localhost:11434"
SCRIPT_DIR  = Path(__file__).parent
PROMPT_FILE = SCRIPT_DIR / "prompts" / "exercise_query.txt"

# Nombres de los modelos personalizados definidos en los Modelfiles
KNOWN_CUSTOM_MODELS = {
    "codellama-custom",
    "deepseek-custom",
    "mistral-custom",
    "phi3-custom",
    "qwen-custom",
}


# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

def get_available_models() -> list[str]:
    """Devuelve los modelos disponibles en la instancia de Ollama."""
    try:
        resp = requests.get(f"{OLLAMA_BASE}/api/tags", timeout=10)
        resp.raise_for_status()
        data = resp.json()
        return [m["name"] for m in data.get("models", [])]
    except requests.exceptions.ConnectionError:
        sys.exit(
            f"No se pudo conectar a Ollama en {OLLAMA_BASE}.\n"
            "Asegurate de que el contenedor del modelo está corriendo."
        )


def query_model(model: str, prompt: str) -> tuple[float, bool, str]:
    """
    Envía el prompt al modelo y devuelve (tiempo_segundos, es_json_válido, respuesta).
    """
    payload = {
        "model": model,
        "prompt": prompt,
        "stream": False,
    }

    start = time.perf_counter()
    resp = requests.post(
        f"{OLLAMA_BASE}/api/generate",
        json=payload,
        timeout=500,   # los modelos 7B pueden tardar
    )
    elapsed = time.perf_counter() - start

    resp.raise_for_status()
    response_text: str = resp.json().get("response", "")

    # Intentar parsear como JSON
    valid = False
    try:
        json.loads(response_text)
        valid = True
    except (json.JSONDecodeError, ValueError):
        valid = False

    return elapsed, valid, response_text


def save_result(model: str, elapsed: float, valid: bool, response: str) -> Path:
    """Guarda el resultado en <model>.json junto al script."""
    output = {
        "time": round(elapsed, 4),
        "valid": valid,
        "response": response,
    }
    out_path = SCRIPT_DIR / f"{model}.json"
    out_path.write_text(json.dumps(output, ensure_ascii=False, indent=2), encoding="utf-8")
    return out_path


# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------

def main() -> None:
    # 1. Leer el prompt
    if not PROMPT_FILE.exists():
        sys.exit(f"No se encontró el archivo de prompt: {PROMPT_FILE}")
    prompt = PROMPT_FILE.read_text(encoding="utf-8").strip()
    print(f"Prompt cargado desde: {PROMPT_FILE}")

    # 2. Detectar modelos disponibles
    available = get_available_models()
    print(f"Modelos disponibles en Ollama: {available}")

    # Filtrar solo los modelos custom de este proyecto
    targets = [m for m in available if m in KNOWN_CUSTOM_MODELS]

    if not targets:
        # Si no hay ninguno de los custom, usar cualquier modelo disponible
        if not available:
            sys.exit("No hay modelos cargados en Ollama.")
        print(
            "No se encontró ningún modelo custom del proyecto. "
            f"Usando el primer modelo disponible: {available[0]}"
        )
        targets = [available[0]]

    # 3. Ejecutar para cada modelo encontrado
    for model in targets:
        print(f"\n[{model}] Enviando prompt...")
        try:
            elapsed, valid, response = query_model(model, prompt)
        except requests.exceptions.Timeout:
            print(f"[{model}] TIMEOUT — el modelo tardó demasiado.")
            continue
        except requests.exceptions.RequestException as exc:
            print(f"[{model}] ERROR en la solicitud: {exc}")
            continue

        out_path = save_result(model, elapsed, valid, response)
        status = "JSON válido ✓" if valid else "JSON inválido ✗"
        print(f"[{model}] Tiempo: {elapsed:.2f}s | {status} | Guardado en: {out_path.name}")


if __name__ == "__main__":
    main()
