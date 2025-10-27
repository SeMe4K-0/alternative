import axios from 'axios'
import type { AxiosInstance, AxiosError } from 'axios'
import { storage } from '../utils'
import type { User, Comet, Observation, Calculation, OrbitParameters, EarthApproach } from '../types'

const API_BASE_URL = (import.meta as any).env?.VITE_API_BASE_URL || '/api'

class ApiService {
  private client: AxiosInstance

  constructor() {
    this.client = axios.create({
      baseURL: API_BASE_URL,
      withCredentials: true,
      headers: {
        'Content-Type': 'application/json',
      },
    })

    this.client.interceptors.request.use((config) => {
      // withCredentials: true автоматически отправляет cookies из браузера
      const sessionId = storage.get<string>('session_id')
      
      // Добавляем session_id в заголовок Cookie, если он есть
      if (sessionId && config.headers) {
        config.headers['Cookie'] = `session_id=${sessionId}`
        console.log('Sending Cookie header with session_id')
      }
      
      return config
    })

    this.client.interceptors.response.use(
      (response) => response,
      (error: AxiosError) => {
        if (error.response?.status === 401) {
          console.log('401 error detected, clearing auth state')
          storage.remove('session_id')
          // НЕ делаем редирект через window.location.href
          // Zustand store сам обработает это через isAuthenticated
        }
        return Promise.reject(error)
      }
    )
  }

  async login(email: string, password: string): Promise<{ user: User; token: string }> {
    try {
      const response = await this.client.post('/auth/login', { email, password })
      
      console.log('Login response headers:', response.headers)
      console.log('All headers keys:', Object.keys(response.headers))
      console.log('Login set-cookie header:', response.headers['set-cookie'])
      
      // Пробуем получить session_id из заголовков
      let sessionId: string | undefined
      const setCookieHeader = response.headers['set-cookie'] as any
      
      console.log('Type of setCookieHeader:', typeof setCookieHeader)
      
      if (setCookieHeader) {
        if (Array.isArray(setCookieHeader)) {
          console.log('setCookieHeader is array, length:', setCookieHeader.length)
          console.log('All cookies:', setCookieHeader)
          const sessionCookie = setCookieHeader.find((c: string) => c.includes('session_id='))
          console.log('Found session cookie:', sessionCookie)
          sessionId = sessionCookie?.split('session_id=')[1]?.split(';')[0]
        } else if (typeof setCookieHeader === 'string') {
          console.log('setCookieHeader is string:', setCookieHeader)
          sessionId = setCookieHeader.split('session_id=')[1]?.split(';')[0]
        }
      } else {
        console.warn('No set-cookie header in response!')
      }
      
      console.log('Extracted session_id:', sessionId)
      console.log('Login response data:', response.data)
      console.log('All response data keys:', Object.keys(response.data))
      
      // Проверяем, какие поля приходят из бэкенда
      console.log('user_id from response:', response.data.user_id)
      console.log('username from response:', response.data.username)
      console.log('avatar_url from response:', response.data.avatar_url)
      console.log('avatar from response:', response.data.avatar)
      
      // Пробуем разные варианты ключей
      const userId = response.data.user_id || response.data.id || response.data.user?.id || 'unknown'
      const userName = response.data.username || response.data.name || response.data.user?.username || email.split('@')[0]
      const userAvatar = response.data.avatar_url || response.data.avatar || response.data.user?.avatar_url || ''
      const userCreatedAt = response.data.created_at || response.data.user?.created_at || new Date().toISOString()
      
      const user: User = {
        id: userId.toString(),
        email,
        name: userName,
        avatar: userAvatar,
        createdAt: userCreatedAt,
      }
      
      console.log('User created from login:', user)
      console.log('User name:', user.name)
      console.log('User avatar:', user.avatar)

      if (sessionId) {
        console.log('Saving session_id to localStorage:', sessionId)
        storage.set('session_id', sessionId)
      } else {
        console.warn('No session_id extracted from login response!')
      }

      return { user, token: sessionId || '' }
    } catch (error: any) {
      console.error('Login error:', error)
      throw new Error(error.response?.data?.error || 'Login failed')
    }
  }

  async register(userData: { name: string; email: string; password: string }): Promise<{ user: User; token: string }> {
    try {
      const response = await this.client.post('/auth/register', {
        email: userData.email,
        password: userData.password,
        username: userData.name,
      })

      const user: User = {
        id: response.data.id.toString(),
        email: response.data.email,
        name: response.data.username || userData.name,
        avatar: response.data.avatar_url,
        createdAt: response.data.created_at,
      }

      // Пробуем получить session_id из заголовков
      let sessionId: string | undefined
      const setCookieHeader = response.headers['set-cookie'] as any
      
      if (setCookieHeader) {
        if (Array.isArray(setCookieHeader)) {
          sessionId = setCookieHeader.find((c: string) => c.includes('session_id='))?.split('session_id=')[1]?.split(';')[0]
        } else if (typeof setCookieHeader === 'string') {
          sessionId = setCookieHeader.split('session_id=')[1]?.split(';')[0]
        }
      }
      
      if (sessionId) {
        storage.set('session_id', sessionId)
      }
      
      console.log('Register response data:', response.data)
      console.log('User created from register:', user)

      return { user, token: sessionId || '' }
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Registration failed')
    }
  }

  async logout(): Promise<void> {
    try {
      await this.client.post('/auth/logout')
    } catch (error) {
      // Ignore errors on logout
      console.error('Logout error:', error)
    }
  }

  async requestPasswordReset(email: string): Promise<void> {
    try {
      await this.client.post('/auth/forgot-password', { email })
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Failed to request password reset')
    }
  }

  async resetPassword(token: string, newPassword: string): Promise<void> {
    try {
      await this.client.post('/auth/reset-password', {
        token,
        new_password: newPassword
      })
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Failed to reset password')
    }
  }

  async getCurrentUser(): Promise<User> {
    try {
      const response = await this.client.get('/profile/me')
      const data = response.data
      
      console.log('getCurrentUser response data:', data)
      
      // Проверяем различные возможные поля для аватара
      const avatarUrl = data.avatar_url || data.avatar || data.image_url
      console.log('Avatar URL from getCurrentUser:', avatarUrl)

      if (!data.id || !data.email) {
        throw new Error('Invalid user data received from server')
      }

      const user: User = {
        id: data.id.toString(),
        email: data.email,
        name: data.username || 'User',
        avatar: avatarUrl,
        createdAt: data.created_at,
      }

      console.log('User loaded from server:', user)
      return user
    } catch (error: any) {
      console.error('Error in getCurrentUser:', error)
      // Не логируем ошибку, если пользователь не авторизован
      if (error.response?.status === 401) {
        console.log('User not authenticated')
      }
      throw new Error(error.response?.data?.error || 'Failed to get user')
    }
  }

  async updateProfile(data: { username?: string; password?: string; avatar?: File }): Promise<User> {
    try {
      console.log('updateProfile called with:', data)
      const formData = new FormData()
      
      if (data.username) {
        console.log('Appending username:', data.username)
        formData.append('username', data.username)
      }
      if (data.password) {
        console.log('Appending password')
        formData.append('password', data.password)
      }
      if (data.avatar) {
        console.log('Appending avatar:', data.avatar.name, data.avatar.type)
        formData.append('avatar', data.avatar)
      }

      console.log('Sending FormData to /profile/me')
      const response = await this.client.put('/profile/me', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
        timeout: 30000, // 30 секунд на загрузку
      })

      console.log('Response received:', response.data)
      const userData = response.data
      
      // Проверяем различные возможные поля для аватара
      console.log('userData.avatar_url:', userData.avatar_url)
      console.log('userData.avatar:', userData.avatar)
      
      const avatarUrl = userData.avatar_url || userData.avatar
      
      const user: User = {
        id: userData.id.toString(),
        email: userData.email,
        name: userData.username,
        avatar: avatarUrl,
        createdAt: userData.created_at,
      }

      console.log('User updated in storage:', user)
      console.log('Avatar URL:', user.avatar)
      return user
    } catch (error: any) {
      console.error('Error in updateProfile:', error)
      console.error('Error response:', error.response?.data)
      throw new Error(error.response?.data?.error || 'Failed to update profile')
    }
  }

  async getComets(): Promise<Comet[]> {
    try {
      const response = await this.client.get('/comets')
      return response.data.map((comet: any) => ({
        id: comet.id.toString(),
        name: comet.name,
        description: comet.description,
        imageUrl: comet.image_url,
        notes: comet.notes || comet.description, // Используем description как notes
        createdAt: comet.created_at,
        updatedAt: comet.updated_at,
      }))
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Failed to get comets')
    }
  }

  async createComet(comet: Omit<Comet, 'id' | 'createdAt' | 'updatedAt'>): Promise<Comet> {
    try {
      // Backend принимает только name и description в JSON формате
      // Поля notes и image не поддерживаются, но отправляем их для будущей поддержки
      const jsonData: any = {
        name: comet.name,
        description: comet.notes || ''
      }

      const response = await this.client.post('/comets', jsonData)

      return {
        id: response.data.id.toString(),
        name: response.data.name,
        description: response.data.description,
        imageUrl: response.data.image_url,
        notes: response.data.notes || response.data.description, // Используем description как notes
        createdAt: response.data.created_at,
        updatedAt: response.data.updated_at,
      }
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Failed to create comet')
    }
  }

  async updateComet(id: string, comet: Partial<Comet>): Promise<Comet> {
    try {
      // Backend принимает только name и description в JSON формате
      const jsonData: any = {
        name: comet.name,
        description: comet.notes || ''
      }

      const response = await this.client.put(`/comets/${id}`, jsonData)

      return {
        id: response.data.id.toString(),
        name: response.data.name,
        description: response.data.description,
        imageUrl: response.data.image_url,
        notes: response.data.description, // Маппим description обратно в notes
        createdAt: response.data.created_at,
        updatedAt: response.data.updated_at,
      }
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Failed to update comet')
    }
  }

  async deleteComet(id: string): Promise<void> {
    try {
      await this.client.delete(`/comets/${id}`)
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Failed to delete comet')
    }
  }

  async getObservations(cometId?: number): Promise<Observation[]> {
    try {
      let response
      if (cometId) {
        response = await this.client.get(`/comets/${cometId}/observations`)
      } else {
        const comets = await this.getComets()
        const allObservations: Observation[] = []
        for (const comet of comets) {
          const obsResponse = await this.client.get(`/comets/${comet.id}/observations`)
          const observations = obsResponse.data.map((obs: any) => ({
            id: obs.id.toString(),
            userId: '1',
            cometName: comet.name,
            observationDate: obs.observed_at,
            coordinates: {
              ra: obs.right_ascension, // Уже в часах, без деления
              dec: obs.declination,
            },
            imageUrl: obs.image_url || '',
            annotatedPoints: [],
            notes: obs.notes,
            createdAt: obs.created_at,
            updatedAt: obs.updated_at,
          }))
          allObservations.push(...observations)
        }
        return allObservations
      }

      return response.data.map((obs: any) => ({
        id: obs.id.toString(),
        userId: '1',
        cometName: '',
        observationDate: obs.observed_at,
        coordinates: {
          ra: obs.right_ascension, // Уже в часах, без деления
          dec: obs.declination,
        },
        imageUrl: obs.image_url || '',
        annotatedPoints: [],
        notes: obs.notes,
        createdAt: obs.created_at,
        updatedAt: obs.updated_at,
      }))
    } catch (error: any) {
      console.error('Failed to get observations:', error)
      return []
    }
  }

  async getObservation(id: string): Promise<Observation> {
    try {
      const response = await this.client.get(`/observations/${id}`)
      const data = response.data

      return {
        id: data.id.toString(),
        userId: '1',
        cometId: '',
        cometName: '',
        observationDate: data.observed_at,
        coordinates: {
          ra: data.right_ascension, // Уже в часах, без деления
          dec: data.declination,
        },
        imageUrl: data.image_url || '',
        annotatedPoints: [],
        notes: data.notes,
        createdAt: data.created_at,
        updatedAt: data.updated_at,
      }
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Failed to get observation')
    }
  }

  async createObservation(observation: Omit<Observation, 'id' | 'createdAt' | 'updatedAt'> & { imageFile?: File }): Promise<Observation> {
    try {
      const formData = new FormData()
      
      // Конвертируем дату в RFC3339 формат
      const observationDate = new Date(observation.observationDate).toISOString()
      formData.append('observed_at', observationDate)
      
      // Отправляем RA в часах без умножения на 15
      formData.append('right_ascension', observation.coordinates.ra.toString())
      
      formData.append('declination', observation.coordinates.dec.toString())

      if (observation.notes) {
        formData.append('notes', observation.notes)
      }

      // Если есть изображение, добавляем его в FormData
      if (observation.imageFile) {
        formData.append('image', observation.imageFile)
      }

      const response = await this.client.post(`/comets/${observation.cometId}/observations`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      })

      return {
        id: response.data.id.toString(),
        userId: '1',
        cometId: observation.cometId,
        cometName: observation.cometName,
        observationDate: response.data.observed_at,
        coordinates: {
          ra: response.data.right_ascension, // Без деления, уже в часах
          dec: response.data.declination,
        },
        imageUrl: response.data.image_url || '',
        annotatedPoints: observation.annotatedPoints || [],
        notes: response.data.notes,
        createdAt: response.data.created_at,
        updatedAt: response.data.updated_at,
      }
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Failed to create observation')
    }
  }

  async updateObservation(id: string, observation: Partial<Observation> & { imageFile?: File }): Promise<Observation> {
    try {
      // Проверяем, есть ли изображение для загрузки
      if ((observation as any).imageFile) {
        // Используем FormData для загрузки изображения
        const formData = new FormData()
        
        if (observation.coordinates) {
          formData.append('right_ascension', observation.coordinates.ra.toString())
          formData.append('declination', observation.coordinates.dec.toString())
        }
        
        if (observation.observationDate) {
          formData.append('observed_at', new Date(observation.observationDate).toISOString())
        }
        
        if (observation.notes !== undefined) {
          formData.append('notes', observation.notes || '')
        }

        formData.append('image', (observation as any).imageFile)

        const response = await this.client.put(`/observations/${id}`, formData, {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        })

        // Нужно получить текущее наблюдение, чтобы сохранить cometName
        const currentObservations = await this.getObservations()
        const currentObservation = currentObservations.find(obs => obs.id === id)
        
        return {
          id: response.data.id.toString(),
          userId: '1',
          cometId: '',
          cometName: currentObservation?.cometName || observation.cometName || '',
          observationDate: response.data.observed_at,
          coordinates: {
            ra: response.data.right_ascension, // Уже в часах, без деления
            dec: response.data.declination,
          },
          imageUrl: response.data.image_url || '',
          annotatedPoints: observation.annotatedPoints || [],
          notes: response.data.notes,
          createdAt: response.data.created_at,
          updatedAt: response.data.updated_at,
        }
      } else {
        // Используем обычный JSON запрос
        const updateData: any = {}
        
        if (observation.coordinates) {
          // Отправляем RA в часах без умножения на 15
          updateData.right_ascension = observation.coordinates.ra
          updateData.declination = observation.coordinates.dec
        }
        
        if (observation.observationDate) {
          // Конвертируем дату в RFC3339 формат
          updateData.observed_at = new Date(observation.observationDate).toISOString()
        }
        
        if (observation.notes !== undefined) {
          updateData.notes = observation.notes
        }

        const response = await this.client.put(`/observations/${id}`, updateData)

        // Нужно получить текущее наблюдение, чтобы сохранить cometName
        const currentObservations = await this.getObservations()
        const currentObservation = currentObservations.find(obs => obs.id === id)

        return {
          id: response.data.id.toString(),
          userId: '1',
          cometId: '',
          cometName: currentObservation?.cometName || observation.cometName || '',
          observationDate: response.data.observed_at,
          coordinates: {
            ra: response.data.right_ascension, // Уже в часах, без деления
            dec: response.data.declination,
          },
          imageUrl: response.data.image_url || '',
          annotatedPoints: observation.annotatedPoints || [],
          notes: response.data.notes,
          createdAt: response.data.created_at,
          updatedAt: response.data.updated_at,
        }
      }
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Failed to update observation')
    }
  }

  async deleteObservation(id: string): Promise<void> {
    try {
      await this.client.delete(`/observations/${id}`)
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Failed to delete observation')
    }
  }

  async getCalculations(cometId?: number): Promise<any[]> {
    try {
      // Загружаем кометы для маппинга comet_id -> название
      const cometsMap = new Map<string, string>()
      const allComets = await this.getComets()
      allComets.forEach(comet => {
        cometsMap.set(comet.id, comet.name)
      })
      
      if (cometId) {
        const response = await this.client.get(`/comets/${cometId}/calculations`)
        return response.data.map((calc: any) => this.mapCalculationResponse(calc, cometsMap))
      } else {
        const allCalculations: any[] = []
        for (const comet of allComets) {
          const calcResponse = await this.client.get(`/comets/${comet.id}/calculations`)
          const calculations = calcResponse.data.map((calc: any) => this.mapCalculationResponse(calc, cometsMap))
          allCalculations.push(...calculations)
        }
        return allCalculations.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())
      }
    } catch (error: any) {
      console.error('Failed to get calculations:', error)
      return []
    }
  }

  async getCalculation(id: string): Promise<Calculation> {
    try {
      const response = await this.client.get(`/calculations/${id}`)
      return this.mapCalculationResponse(response.data)
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Failed to get calculation')
    }
  }

  // Этап 1: Расчет параметров орбиты
  async createCalculation(observationIds: string[], cometId?: number): Promise<Calculation> {
    try {
      let targetCometId = cometId

      if (!targetCometId) {
        const comets = await this.getComets()
        if (comets.length === 0) {
          throw new Error('No comets found. Please create a comet first.')
        }
        targetCometId = parseInt(comets[0].id)
      }

      const response = await this.client.post(`/comets/${targetCometId}/calculations`, {
        observation_ids: observationIds.map((id) => parseInt(id)),
      })

      const data = response.data
      
      // Возвращаем только параметры орбиты (без сближения)
      return {
        id: data.id.toString(),
        userId: '1',
        observationIds: observationIds,
        orbitParameters: {
          semiMajorAxis: data.semi_major_axis,
          eccentricity: data.eccentricity,
          inclination: data.inclination,
          longitudeOfAscendingNode: data.lon_ascending_node,
          argumentOfPeriapsis: data.arg_periapsis,
          meanAnomaly: data.mean_anomaly || 0,
          period: Math.sqrt(Math.pow(data.semi_major_axis, 3)), // Расчет периода
        },
        earthApproach: undefined, // Пока нет сближения
        status: 'pending',
        createdAt: data.created_at || data.calculated_at,
        completedAt: data.calculated_at,
      }
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Failed to create calculation')
    }
  }
  
  // Этап 2: Расчет сближения с Землей
  async calculateApproach(calculationId: string, orbitParams: OrbitParameters): Promise<EarthApproach> {
    try {
      const response = await this.client.post(`/calculations/${calculationId}/approach`, {
        semi_major_axis: orbitParams.semiMajorAxis,
        eccentricity: orbitParams.eccentricity,
        inclination: orbitParams.inclination,
        lon_ascending_node: orbitParams.longitudeOfAscendingNode,
        arg_periapsis: orbitParams.argumentOfPeriapsis,
        time_perihelion: new Date().toISOString(), // Временное значение
      })

      const data = response.data
      return {
        date: data.approach_date,
        distance: data.distance_au,
        velocity: 0, // Временно
      }
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Failed to calculate approach')
    }
  }

  private mapCalculationResponse(data: any, cometsMap?: Map<string, string>): any {
    // Логируем структуру данных для отладки
    console.log('=== Calculation data structure ===')
    console.log('All keys:', Object.keys(data))
    console.log('time_perihelion:', data.time_perihelion)
    console.log('timePerihelion:', data.timePerihelion)
    console.log('approximate_time:', data.approximate_time)
    console.log('Full data:', JSON.stringify(data, null, 2))
    console.log('===================================')
    
    // Получаем название кометы из первого наблюдения
    let cometName = null
    if (data.observations && data.observations.length > 0) {
      const obs = data.observations[0]
      cometName = obs.comet?.name || obs.comet_name || obs.Comet?.name || null
    }
    
    // Если не нашли в наблюдении, используем comet_id
    if (!cometName && data.comet_id && cometsMap) {
      cometName = cometsMap.get(data.comet_id.toString()) || null
    }
    
    console.log('Extracted comet name:', cometName)
    
    return {
      id: data.id.toString(),
      userId: '1',
      observationIds: data.observations?.map((obs: any) => obs.id.toString()) || [],
      cometName: cometName, // Добавляем название кометы
      orbitParameters: {
        semiMajorAxis: data.semi_major_axis,
        eccentricity: data.eccentricity,
        inclination: data.inclination,
        longitudeOfAscendingNode: data.lon_ascending_node,
        argumentOfPeriapsis: data.arg_periapsis,
        meanAnomaly: 0,
        period: 0,
        timePerihelion: data.time_perihelion || data.timePerihelion || data.approximate_time,
      },
      earthApproach: data.approach_date
        ? {
            date: data.approach_date,
            distance: data.distance_au || data.distance_km || 0,
            velocity: 0,
          }
        : undefined,
      status: 'completed',
      createdAt: data.created_at || data.calculated_at,
      completedAt: data.calculated_at,
    }
  }

  async deleteCalculation(id: string): Promise<void> {
    try {
      await this.client.delete(`/calculations/${id}`)
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Failed to delete calculation')
    }
  }

  // Метод uploadImage удален, так как загрузка изображений происходит через createObservation
}

export const apiService = new ApiService()