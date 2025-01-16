# Autorzy:
# Jan Wolski s26435
# Marcin Topolniak s25672

# Eye Tracker:
# - Program wykorzystuje kamerę do śledzenia ruchu oczu użytkownika
# - Na podstawie analizy punktów charakterystycznych twrzy oblicza, gdzie na ekranie patrzy użytkownik
# - Generuje heatmape, która pokazuje obszary najczęściej obserwowane przez użytkownika

# Link do nagrania wideo znajduje się w pliku read.me w folderze z zadaniem

import cv2
import mediapipe as mp
import numpy as np
import time

# Funkcja, która oblicza wektor spojrzenia i położenie oczu użytkownika na podstawie punktów charakterystycznych wykrytych przez MediaPipe
def calculate_gaze_direction(landmarks, image_width, image_height):
    # Kluczowe punkty wokół oka (górna i dolna powieka, rogi oka)
    left_eye_points = [landmarks[33], landmarks[133], landmarks[159], landmarks[145]]  # Lewe oko
    right_eye_points = [landmarks[362], landmarks[263], landmarks[386], landmarks[374]]  # Prawe oko

    # Obliczanie środków oczu jako średniej punktów chrakterystycznych
    left_eye_center = np.mean(np.array([[p.x * image_width, p.y * image_height] for p in left_eye_points]), axis=0)
    right_eye_center = np.mean(np.array([[p.x * image_width, p.y * image_height] for p in right_eye_points]), axis=0)

    # Pozycje źrenic (około środka oka)
    left_pupil = landmarks[468]  # Źrenica lewego oka
    right_pupil = landmarks[473]  # Źrenica prawego oka

    # Przekształcenie współrzędnych źrenic z normalizowanych na piksele obrazu kamery
    left_pupil_position = np.array([left_pupil.x * image_width, left_pupil.y * image_height])
    right_pupil_position = np.array([right_pupil.x * image_width, right_pupil.y * image_height])

    # Wektor spojrzenia w każdym oku
    left_gaze_vector = left_pupil_position - left_eye_center
    right_gaze_vector = right_pupil_position - right_eye_center

    # Średni wektor spojrzenia
    gaze_vector = (left_gaze_vector + right_gaze_vector) / 2

    # Przeskalowanie i przesunięcie na środek ekranu
    # Odwrócenie osi Y dla poprawnej orientacji
    gaze_vector[1] = -gaze_vector[1]
    
    return gaze_vector, (left_eye_center + right_eye_center) / 2


# Funkcja do zarządzania kamerą, przetwarzania obrazów, śledzenia spojrzenia i generowania heatmapy
def main():
    # Inicjalizacja modelu MediaPipe FaceMesh
    mp_face_mesh = mp.solutions.face_mesh
    face_mesh = mp_face_mesh.FaceMesh(refine_landmarks=True, min_detection_confidence=0.5, min_tracking_confidence=0.5)

    # Wczytaj obraz reklamy
    advertisement = cv2.imread('reklama.png')
    if advertisement is None:
        print("Nie można wczytać pliku reklama.png.")
        return
    
    # Uruchomienie kamery
    cap = cv2.VideoCapture(0)
    if not cap.isOpened():
        print("Nie można otworzyć kamery.")
        return

    # Pobranie pierwszej klatki, aby ustalić rozdzielczość obrazu
    ret, frame = cap.read()
    if not ret:
        print("Nie można odczytać klatki z kamery.")
        return

    frame_height, frame_width, _ = frame.shape

    # Dopasowanie reklamy i mapy ciepła do rozmiarów obrazu z kamery
    resized_advertisement = cv2.resize(advertisement, (frame_width, frame_height))
    heatmap = np.zeros((frame_height, frame_width), dtype=np.float32)

    # Zapisanie czasu początkowego
    start_time = time.time()

    while cap.isOpened():
        # Odczyt kolejnej klatki z kamery
        ret, frame = cap.read()
        if not ret:
            break
        # Konwersja koloru na RGB dla MediaPipe
        frame_rgb = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)
        results = face_mesh.process(frame_rgb)

        if results.multi_face_landmarks:
            # Iteracja po wykrytej twarzy
            for face_landmarks in results.multi_face_landmarks:
                # Obliczanie wektora spojrzenia i punktu początkowego
                gaze_vector, gaze_origin = calculate_gaze_direction(face_landmarks.landmark, frame_width, frame_height)

                # Obliczenie punktu spojrzenia na obrazie reklamy
                gaze_point = gaze_origin + gaze_vector * 100  # Skalowanie wektora dla wizualizacji
                x, y = int(gaze_point[0]), int(gaze_point[1])

                # Aktualizacja heatmapy
                if 0 <= x < frame_width and 0 <= y < frame_height:
                    cv2.circle(heatmap, (x, y), 15, (255), -1)  # Większe, bardziej widoczne punkty(r= 15px)

                # Nanoszenie punktu spojrzenia na obraz z kamery
                cv2.circle(frame, (int(gaze_origin[0]), int(gaze_origin[1])), 5, (0, 255, 0), -1)
                cv2.line(frame, (int(gaze_origin[0]), int(gaze_origin[1])), (x, y), (0, 0, 255), 2)

        # Wyświetlenie obrazu reklamy w fullscreenie
        cv2.namedWindow('Advertisement', cv2.WND_PROP_FULLSCREEN)
        cv2.setWindowProperty('Advertisement', cv2.WND_PROP_FULLSCREEN, cv2.WINDOW_FULLSCREEN)
        cv2.imshow('Advertisement', resized_advertisement)

        # Wyświetlenie obrazu z kamery z zaznaczonym punktem spojrzenia
        cv2.imshow('Camera View', frame)

        # Wyłącz po 6 sekundach
        if time.time() - start_time > 6:
            break
        
        # Wyłącz po naciśnięciu klawisza "q"
        if cv2.waitKey(1) & 0xFF == ord('q'):
            break
    # Zwolnienie zasobów kamery i zamknięcie okien
    cap.release()
    cv2.destroyAllWindows()

    # Normalizacja heatmapy
    heatmap = cv2.normalize(heatmap, None, 0, 255, cv2.NORM_MINMAX).astype(np.uint8)

    # Konwersja heatmapy na obraz kolorowy
    heatmap_color = cv2.applyColorMap(heatmap, cv2.COLORMAP_JET)
   
    # Dopasowanie heatmapy do rozmiaru reklamy
    heatmap_resized = cv2.resize(heatmap_color, (frame_width, frame_height))

    # Nałożenie heatmapy na obraz reklamy z jaskrawo czerwonym kolorem
    combined = cv2.addWeighted(resized_advertisement, 0.6, heatmap_resized, 0.4, 0)

    # Wyświetlenie wyniku mapy ciepła w osobnym oknie
    cv2.imshow('Heatmap on Advertisement', combined)
    cv2.waitKey(0)
    cv2.destroyAllWindows()

if __name__ == "__main__":
    main()
