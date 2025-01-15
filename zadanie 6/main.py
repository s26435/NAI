import cv2
import mediapipe as mp
import numpy as np
import time

def calculate_gaze_direction(landmarks, image_width, image_height):
    # Kluczowe punkty wokół oka
    left_eye_points = [landmarks[33], landmarks[133], landmarks[159], landmarks[145]]  # Lewe oko
    right_eye_points = [landmarks[362], landmarks[263], landmarks[386], landmarks[374]]  # Prawe oko

    # Środek oczu
    left_eye_center = np.mean(np.array([[p.x * image_width, p.y * image_height] for p in left_eye_points]), axis=0)
    right_eye_center = np.mean(np.array([[p.x * image_width, p.y * image_height] for p in right_eye_points]), axis=0)

    # Pozycje źrenic (około środka oka)
    left_pupil = landmarks[468]  # Źrenica lewego oka
    right_pupil = landmarks[473]  # Źrenica prawego oka

    left_pupil_position = np.array([left_pupil.x * image_width, left_pupil.y * image_height])
    right_pupil_position = np.array([right_pupil.x * image_width, right_pupil.y * image_height])

    # Wektor spojrzenia w każdym oku
    left_gaze_vector = left_pupil_position - left_eye_center
    right_gaze_vector = right_pupil_position - right_eye_center

    # Średni wektor spojrzenia
    gaze_vector = (left_gaze_vector + right_gaze_vector) / 2

    # Przeskalowanie i przesunięcie na środek ekranu
    gaze_vector[1] = -gaze_vector[1]  # Odwrócenie osi Y dla poprawnej orientacji
    return gaze_vector, (left_eye_center + right_eye_center) / 2

def main():
    mp_face_mesh = mp.solutions.face_mesh
    face_mesh = mp_face_mesh.FaceMesh(refine_landmarks=True, min_detection_confidence=0.5, min_tracking_confidence=0.5)

    # Wczytaj obraz reklamy
    advertisement = cv2.imread('reklama.png')
    if advertisement is None:
        print("Nie można wczytać pliku reklama.png.")
        return

    cap = cv2.VideoCapture(0)

    if not cap.isOpened():
        print("Nie można otworzyć kamery.")
        return

    # Pobierz rozmiar ramki kamery
    ret, frame = cap.read()
    if not ret:
        print("Nie można odczytać klatki z kamery.")
        return

    frame_height, frame_width, _ = frame.shape

    # Dopasowanie reklamy i mapy ciepła do rozmiaru ekranu
    resized_advertisement = cv2.resize(advertisement, (frame_width, frame_height))
    heatmap = np.zeros((frame_height, frame_width), dtype=np.float32)

    start_time = time.time()

    while cap.isOpened():
        ret, frame = cap.read()
        if not ret:
            break

        frame_rgb = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)
        results = face_mesh.process(frame_rgb)

        if results.multi_face_landmarks:
            for face_landmarks in results.multi_face_landmarks:
                gaze_vector, gaze_origin = calculate_gaze_direction(face_landmarks.landmark, frame_width, frame_height)

                # Obliczenie punktu spojrzenia na obrazie reklamy
                gaze_point = gaze_origin + gaze_vector * 100  # Skalowanie wektora dla wizualizacji
                x, y = int(gaze_point[0]), int(gaze_point[1])

                if 0 <= x < frame_width and 0 <= y < frame_height:
                    cv2.circle(heatmap, (x, y), 15, (255), -1)  # Większe, bardziej widoczne punkty

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

        if cv2.waitKey(1) & 0xFF == ord('q'):
            break

    cap.release()
    cv2.destroyAllWindows()

    # Normalizacja mapy częstości
    heatmap = cv2.normalize(heatmap, None, 0, 255, cv2.NORM_MINMAX).astype(np.uint8)

    # Dopasowanie mapy ciepła do rozmiaru reklamy
    heatmap_color = cv2.applyColorMap(heatmap, cv2.COLORMAP_JET)
    heatmap_resized = cv2.resize(heatmap_color, (frame_width, frame_height))

    # Nałożenie mapy częstości na obraz reklamy z jaskrawo czerwonym kolorem
    combined = cv2.addWeighted(resized_advertisement, 0.6, heatmap_resized, 0.4, 0)

    # Wyświetlenie wyniku mapy ciepła w osobnym oknie
    cv2.imshow('Heatmap on Advertisement', combined)
    cv2.waitKey(0)
    cv2.destroyAllWindows()

if __name__ == "__main__":
    main()
