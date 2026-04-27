type EventName = 'session:expired'
type Listener = () => void

class EventBus {
  private listeners = new Map<EventName, Set<Listener>>()

  on(event: EventName, listener: Listener): () => void {
    if (!this.listeners.has(event)) this.listeners.set(event, new Set())
    this.listeners.get(event)!.add(listener)
    return () => this.off(event, listener)
  }

  off(event: EventName, listener: Listener): void {
    this.listeners.get(event)?.delete(listener)
  }

  emit(event: EventName): void {
    this.listeners.get(event)?.forEach((l) => l())
  }
}

export const eventBus = new EventBus()
