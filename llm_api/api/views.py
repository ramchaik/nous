from django.http import JsonResponse
from django.views.decorators.csrf import csrf_exempt
from django.views.decorators.http import require_http_methods
import json
from .main import LangChainModel

@csrf_exempt
@require_http_methods(["POST"])
def predict(request):
    try:
        data = json.loads(request.body)
        question = data.get('question')
        
        if not question:
            return JsonResponse({'error': 'Question is required'}, status=400)
        
        model = LangChainModel()
        result = model.predict(question)
        
        return JsonResponse(result)
    except Exception as e:
        return JsonResponse({'error': str(e)}, status=500)