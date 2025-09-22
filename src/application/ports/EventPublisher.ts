export interface EventPublisher {
  publish<T>(topic: string, payload: T): Promise<void>;
}
