package entity

type World struct {
    mutex sync.RWMutex
    Players map[int64]*Player
    //Enemys map[int64]*Enemy
    //Items map[string]*WorldObject 

    TickCounter int64          
    Config      *WorldConfig 
    IsRunning   atomic.Bool 
}