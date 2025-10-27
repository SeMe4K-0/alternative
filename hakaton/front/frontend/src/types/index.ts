export interface User {
  id: string
  email: string
  name: string
  avatar?: string
  apiKey?: string
  createdAt: string
}

export interface Point {
  x: number
  y: number
  label?: string
}

export interface Comet {
  id: string
  name: string
  description?: string
  imageUrl?: string
  notes?: string
  createdAt: string
  updatedAt: string
}

export interface Observation {
  id: string
  userId: string
  cometId: string
  cometName: string
  observationDate: string
  coordinates: {
    ra: number // Right Ascension in hours
    dec: number // Declination in degrees
  }
  imageUrl: string
  annotatedPoints: Point[]
  notes?: string
  createdAt: string
  updatedAt: string
}

export interface Calculation {
  id: string
  userId: string
  observationIds: string[]
  cometName?: string
  orbitParameters: OrbitParameters
  earthApproach?: EarthApproach
  status: 'pending' | 'completed' | 'failed'
  createdAt: string
  completedAt?: string
}

export interface OrbitParameters {
  semiMajorAxis: number // AU
  eccentricity: number
  inclination: number // degrees
  longitudeOfAscendingNode: number // degrees
  argumentOfPeriapsis: number // degrees
  meanAnomaly: number // degrees
  period: number // years
  timePerihelion?: string // ISO date string
}

export interface EarthApproach {
  date: string
  distance: number // AU
  velocity: number // km/s
}

export interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  isLoading: boolean
}

export interface ObservationState {
  observations: Observation[]
  currentObservation: Observation | null
  isLoading: boolean
  error: string | null
}

export interface CalculationState {
  calculations: Calculation[]
  currentCalculation: Calculation | null
  isLoading: boolean
  error: string | null
}

export interface UIState {
  theme: 'light' | 'dark'
  sidebarOpen: boolean
  notifications: Notification[]
}

export interface Notification {
  id: string
  type: 'success' | 'error' | 'warning' | 'info'
  title: string
  message: string
  duration?: number
}

export interface PasswordResetRequest {
  email: string
}

export interface PasswordResetConfirmRequest {
  token: string
  newPassword: string
}
