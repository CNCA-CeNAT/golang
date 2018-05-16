#!/opt/intel/intelpython35/bin/python3

#PBS -N golang-benchmark
#PBS -q phi-n1h72
#PBS -l nodes=1:ppn=1
#PBS -l walltime=20:00:00

import scalability
import modules
import os

cores = [1,2,4,8,16,32,64]

def montecarlo_go():
    path = os.path.expanduser('~/golang-benchmarks/')
    P = lambda c, s: (path + 'montecarlo_pi %d %d') % (c, s)
    setup = lambda: modules.run('go build -o %smontecarlo_pi %smontecarlo_pi.go' % (path, path))
    montecarlo = scalability.Executable(P, setup, 'montecarlo go')
    G_strong = lambda c: 600000000
    G_weak   = scalability.poli_growth(100000000, 1)
    scalability.scaling_full_test(montecarlo, G_strong, G_weak, cores, 10)

scalability.scalability_main(globals())
