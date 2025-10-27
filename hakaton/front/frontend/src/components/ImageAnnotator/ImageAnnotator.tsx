import React, { useRef, useEffect, useState, useCallback } from 'react'
import { Stage, Layer, Image as KonvaImage, Circle, Text } from 'react-konva'
import type { Point } from '../../types/index'

interface ImageAnnotatorProps {
  imageUrl: string
  points: Point[]
  onPointsChange: (points: Point[]) => void
  width?: number
  height?: number
  disabled?: boolean
}

const ImageAnnotator: React.FC<ImageAnnotatorProps> = ({
  imageUrl,
  points,
  onPointsChange,
  width = 600,
  height = 400,
  disabled = false
}) => {
  const stageRef = useRef<any>(null)
  const imageRef = useRef<any>(null)
  const [image, setImage] = useState<HTMLImageElement | null>(null)
  const [stageSize, setStageSize] = useState({ width, height })
  const [imageSize, setImageSize] = useState({ width: 0, height: 0 })

  // Загрузка изображения
  useEffect(() => {
    const img = new Image()
    img.crossOrigin = 'anonymous'
    img.onload = () => {
      setImage(img)
      
      // Вычисляем размеры для отображения с сохранением пропорций
      const aspectRatio = img.width / img.height
      let displayWidth = width
      let displayHeight = height
      
      if (aspectRatio > width / height) {
        displayHeight = width / aspectRatio
      } else {
        displayWidth = height * aspectRatio
      }
      
      setImageSize({ width: displayWidth, height: displayHeight })
      setStageSize({ width: displayWidth, height: displayHeight })
    }
    img.src = imageUrl
  }, [imageUrl, width, height])

  // Обработка клика по изображению
  const handleStageClick = useCallback((e: any) => {
    if (disabled || !image) return

    const stage = stageRef.current
    const pointerPosition = stage.getPointerPosition()
    
    if (!pointerPosition) return

    // Проверяем, что клик был по изображению
    const x = pointerPosition.x
    const y = pointerPosition.y
    
    if (x < 0 || y < 0 || x > imageSize.width || y > imageSize.height) {
      return
    }

    // Добавляем новую точку
    const newPoint: Point = {
      x,
      y,
      label: `Точка ${points.length + 1}`
    }

    onPointsChange([...points, newPoint])
  }, [disabled, image, points, onPointsChange, imageSize])

  // Удаление точки
  const handlePointClick = useCallback((index: number) => {
    if (disabled) return
    
    const newPoints = points.filter((_, i) => i !== index)
    onPointsChange(newPoints)
  }, [disabled, points, onPointsChange])

  // Очистка всех точек
  const clearPoints = useCallback(() => {
    if (disabled) return
    onPointsChange([])
  }, [disabled, onPointsChange])

  if (!image) {
    return (
      <div className="flex items-center justify-center bg-gray-100 rounded-lg" style={{ width, height }}>
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600 mx-auto mb-2"></div>
          <p className="text-gray-600">Загрузка изображения...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h3 className="text-lg font-medium text-gray-900">
          Аннотация изображения
        </h3>
        <div className="space-x-2">
          <button
            onClick={clearPoints}
            disabled={disabled || points.length === 0}
            className="px-3 py-1 text-sm bg-red-100 text-red-700 rounded-md hover:bg-red-200 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Очистить все
          </button>
          <span className="text-sm text-gray-500">
            Точек: {points.length}
          </span>
        </div>
      </div>

      <div className="border border-gray-300 rounded-lg overflow-hidden bg-gray-50">
        <Stage
          ref={stageRef}
          width={stageSize.width}
          height={stageSize.height}
          onClick={handleStageClick}
          style={{ cursor: disabled ? 'default' : 'crosshair' }}
        >
          <Layer>
            <KonvaImage
              ref={imageRef}
              image={image}
              width={imageSize.width}
              height={imageSize.height}
            />
            
            {/* Отображение точек */}
            {points.map((point, index) => (
              <React.Fragment key={index}>
                <Circle
                  x={point.x}
                  y={point.y}
                  radius={8}
                  fill="#ef4444"
                  stroke="#ffffff"
                  strokeWidth={2}
                  onClick={() => handlePointClick(index)}
                  style={{ cursor: disabled ? 'default' : 'pointer' }}
                />
                <Text
                  x={point.x + 12}
                  y={point.y - 8}
                  text={point.label || `${index + 1}`}
                  fontSize={12}
                  fill="#374151"
                  fontStyle="bold"
                />
              </React.Fragment>
            ))}
          </Layer>
        </Stage>
      </div>

      <div className="text-sm text-gray-600">
        <p>• Кликните по изображению, чтобы добавить точку</p>
        <p>• Кликните по красной точке, чтобы удалить её</p>
        <p>• Минимум 3 точки необходимо для расчёта орбиты</p>
      </div>
    </div>
  )
}

export default ImageAnnotator
